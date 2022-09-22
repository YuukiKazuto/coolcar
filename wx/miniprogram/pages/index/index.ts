import { CarService } from "../../service/car"
import { ProfileService } from "../../service/profile"
import { rental } from "../../service/proto_gen/rental/rental_pb"
import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"

interface Marker {
  iconPath: string
  id: number
  latitude: number
  longitude: number
  width: number
  height: number
}
// const defaultAvatar = 'https://mapapi.qq.com/web/mapComponents/componentsTest/zyTest/static/model_taxi.png?a=1'
const defaultAvatar = '/resources/car.png'

const initialLat = 24.5
const initialLng = 118.11
Page({
  socket: undefined as WechatMiniprogram.SocketTask | undefined,
  isPageShowing: false,
  data: {
    setting: {
      skew: 0,
      rotate: 0,
      showLocation: true,
      showScale: true,
      subKey: '',
      layerStyle: -1,
      enableZoom: true,
      enableScroll: true,
      enableRotate: false,
      showCompass: false,
      enable3D: false,
      enableOverlooking: false,
      enableSatellite: false,
      enableTraffic: false,
    },
    location: {
      latitude: initialLat,
      longitude: initialLng,
    },
    is3D: false,
    isOverlooking: false,
    scale: 12,
    markers: [] as Marker[]
  },

  onLoad() {
  },

  onShow() {
    this.isPageShowing = true
    if (!this.socket) {
      this.setData({
        markers: []
      }, () => {
        this.setupCarPosUpdater()
      })
    }
  },

  onHide() {
    this.isPageShowing = false
    if (this.socket) {
      this.socket.close({
        success: () => {
          this.socket = undefined
        }
      })
    }
  },

  async onScanTap() {
    const tripsRes = await TripService.getTrips(rental.v1.TripStatus.IN_PROGRESS)
    if ((tripsRes.trips?.length || 0) > 0) {
      await this.selectComponent('#tripModal').showModal()
      wx.navigateTo({
        url: routing.driving({
          trip_id: tripsRes.trips![0].id!
        })
      })
      return
    }
    wx.scanCode({
      success: async (res) => {
        const carID = res.result
        const lockURL = routing.lock({
          car_id: carID
        })
        const prof = await ProfileService.getProfile()
        if (prof.identityStatus === rental.v1.IdentityStatus.VERIFIED) {
          wx.navigateTo({
            url: lockURL,
          })
        } else {
          await this.selectComponent('#licModal').showModal()
          wx.navigateTo({
            url: routing.register({
              redirectURL: lockURL,
            })
          })
        }
      },
      fail: console.error
    })
  },

  onMyLocationTap() {
    wx.getLocation({
      type: 'gcj02',
      success: res => {
        this.setData({
          location: {
            latitude: res.latitude,
            longitude: res.longitude
          }
        })
      },
      fail: () => {
        wx.showToast({
          icon: 'none',
          title: '请前往设置页授权'
        })
      }
    })
  },

  onMyTripsTap() {
    wx.navigateTo({
      url: routing.mytrips(),
    })
  },

  setupCarPosUpdater(){
    const map = wx.createMapContext("map")
    const markersByCarID = new Map<string, Marker>()
    let translationInProgress = false
    const endTranslation = () => {
      translationInProgress = false
    }
    this.socket = CarService.subscribe(car => {
      if (!car.id || translationInProgress || !this.isPageShowing) {
        console.log('dropped')
        return
      }
      const marker = markersByCarID.get(car.id)
      if (!marker) {
        // Insert new marker.
        const newMarker: Marker = {
          id: this.data.markers.length,
          iconPath: car.car?.driver?.avatarUrl || defaultAvatar,
          latitude: car.car?.position?.latitude || initialLat,
          longitude: car.car?.position?.longitude || initialLng,
          height: 40,
          width: 40,
        }
        markersByCarID.set(car.id, newMarker)
        this.data.markers.push(newMarker)
        translationInProgress = true
        this.setData({
          markers: this.data.markers,
        }, endTranslation)
        return
      }

      const newAvatar = car.car?.driver?.avatarUrl || defaultAvatar
      const newLat = car.car?.position?.latitude || initialLat
      const newLng = car.car?.position?.longitude || initialLng
      if (marker.iconPath !== newAvatar) {
        // Change iconPath and possibly position.
        marker.iconPath = newAvatar
        marker.latitude = newLat
        marker.longitude = newLng
        translationInProgress = true
        this.setData({
          markers: this.data.markers,
        }, endTranslation)
        return
      }

      if (marker.latitude !== newLat || marker.longitude !== newLng) {
        // Move marker.
        translationInProgress = true
        map.translateMarker({
          markerId: marker.id,
          destination: {
            latitude: newLat,
            longitude: newLng,
          },
          autoRotate: false,
          rotate: 0,
          duration: 800,
          animationEnd: endTranslation,
        })
      }
    })
  },
})
