import { CarService } from "../../service/car"
import { car } from "../../service/proto_gen/car/car_pb"
import { rental } from "../../service/proto_gen/rental/rental_pb"
import { TripService } from "../../service/trip"
import { routing } from "../../utils/routing"

// pages/lock/lock.ts
Page({
  carID: '',
  carRefresher: 0,
  data: {
    avatarURL: '',
    shareLocation: false
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad(opt: Record<'car_id', string>) {
    const o: routing.LockOpts = opt
    this.carID = o.car_id
  },

  getUserProfile() {
    wx.getUserProfile({
      desc: '允许其他用户看到我此次行程头像',
      success: (res) => {
        this.setData({
          avatarURL: res.userInfo.avatarUrl,
          shareLocation: true
        })
      },
      fail: console.log
    })
  },
  onShareLocation(e: any) {
    this.data.shareLocation = e.detail.value
    if (!this.data.shareLocation) {
      this.setData({
        avatarURL: '',
      })
    }
  },
  onUnlockTap() {
    wx.getLocation({
      type: 'gcj02',
      success: async loc => {
        if (!this.carID) {
          console.error('no carID specified');
          return
        }
        let trip: rental.v1.ITripEntity
        try {
          trip = await TripService.createTrip({
            start: loc,
            carId: this.carID,
            avatarUrl: this.data.shareLocation ? this.data.avatarURL : '',
          })

          if (!trip.id) {
            console.error('no tripID in response', trip)
            return
          }
        } catch (error) {
          wx.showToast({
            title: '创建行程失败',
            icon: 'error'
          })
        }

        wx.showLoading({
          title: '开锁中',
          mask: true,
        })

        this.carRefresher = setInterval(async () => {
          const c = await CarService.getCar(this.carID)
          if (c.status === car.v1.CarStatus.UNLOCKED) {
            this.clearCarRefresher()
            wx.redirectTo({
              url: routing.driving({
                trip_id: trip.id!,
              }),
              complete: () => {
                wx.hideLoading()
              }
            })
          }
        }, 2000)
      },
      fail: () => {
        wx.showToast({
          icon: 'none',
          title: '请前往设置页授权位置信息'
        })
      }
    })
  },
  onUnload() {
    this.clearCarRefresher()
    wx.hideLoading()
  },

  clearCarRefresher() {
    if (this.carRefresher) {
      clearInterval(this.carRefresher)
      this.carRefresher = 0
    }
  },
})