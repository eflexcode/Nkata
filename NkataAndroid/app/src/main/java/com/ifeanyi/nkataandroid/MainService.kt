package com.ifeanyi.nkataandroid

import android.app.PendingIntent
import android.app.Service
import android.content.Intent
import android.os.IBinder
import androidx.core.app.NotificationCompat
import androidx.lifecycle.compose.rememberLifecycleOwner
import androidx.room.Room
import com.ifeanyi.nkataandroid.logic.Util
import com.ifeanyi.nkataandroid.logic.database.room.NkataDatabase
import com.ifeanyi.nkataandroid.logic.network.okhttp.WsListener
import com.ifeanyi.nkataandroid.logic.network.retrofit.apiservice.NkataRetrofitClient
import com.ifeanyi.nkataandroid.logic.repository.LogicRepository
import com.ifeanyi.nkataandroid.logic.repository.impl.LogicRepositoryImpl
import com.ifeanyi.nkataandroid.ui.viewmodel.NkataViewModel
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.WebSocket
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

class MainService : Service() {

    lateinit var appDb: NkataDatabase
    lateinit var webSocket: WebSocket
    lateinit var okHttpClient: OkHttpClient

    override fun onCreate() {
        super.onCreate()

        val repo = LogicRepositoryImpl(appDb.dao())

        val v = NkataViewModel(repo)

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

        val intent = Intent(this, MainActivity::class.java)
        val pendingIntent = PendingIntent.getActivity(
            this,
            0,
            intent,
            PendingIntent.FLAG_IMMUTABLE // Required for Android 12+
        )

        val notificationCompat = NotificationCompat.Builder(this, Util.CHANNEL_ID)
            .setContentTitle("Chat Service Active")
            .setContentText("Connected and waiting for messages...")
            .setSmallIcon(android.R.drawable.stat_notify_chat)
            .setContentIntent(pendingIntent)
            .setOngoing(true)
            .setCategory(NotificationCompat.CATEGORY_SERVICE)
            .build()

        startForeground(Util.NOTIFICATION_ID_FOREGROUND, notificationCompat)
        return START_STICKY
    }

    override fun onDestroy() {
        super.onDestroy()
        webSocket.close(1000, "Service killed")
    }

    override fun onBind(intent: Intent): IBinder? {
        return null
    }
}