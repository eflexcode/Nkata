package com.ifeanyi.nkataandroid

import android.app.Service
import android.content.Intent
import android.os.IBinder
import androidx.room.Room
import com.ifeanyi.nkataandroid.logic.Util
import com.ifeanyi.nkataandroid.logic.database.room.NkataDatabase
import com.ifeanyi.nkataandroid.logic.network.okhttp.WsListener
import okhttp3.OkHttp
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.WebSocket

class MainService : Service() {

    lateinit var appDb: NkataDatabase
    lateinit var webSocket: WebSocket
    lateinit var okHttpClient: OkHttpClient

    override fun onCreate() {
        super.onCreate()

        appDb = Room.databaseBuilder(
            application,
            NkataDatabase::class.java,
            Util.DatabaseName
        ).build()

        okHttpClient = OkHttpClient()

        val f = "http://localhots:5557/v1/general-authenticated/ws/{user_id}"
        val request = Request.Builder()
        request.addHeader("Authorization", "token")// get token from db
        request.url(f)

        val wsListener = WsListener()
        webSocket = okHttpClient.newWebSocket(request.build(), wsListener)

    //get both edittext and img data from viewmodel for webSocket.send()

    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        startForeground(NOTIFICATION_ID, createNotification())
       return START_STICKY
    }

    override fun onDestroy() {
        super.onDestroy()
        webSocket.close(1000,"Service killed")
    }
    override fun onBind(intent: Intent): IBinder {
        return null
    }
}