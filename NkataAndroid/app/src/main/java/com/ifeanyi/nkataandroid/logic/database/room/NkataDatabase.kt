package com.ifeanyi.nkataandroid.logic.database.room

import androidx.room.Database
import androidx.room.RoomDatabase
import com.ifeanyi.nkataandroid.logic.database.room.dao.NkataDatabaseDao
import com.ifeanyi.nkataandroid.logic.database.room.model.Chat
import com.ifeanyi.nkataandroid.logic.database.room.model.Friendship
import com.ifeanyi.nkataandroid.logic.database.room.model.Message
import com.ifeanyi.nkataandroid.logic.database.room.model.Notification
import com.ifeanyi.nkataandroid.logic.database.room.model.Profile
import com.ifeanyi.nkataandroid.logic.database.room.model.Token
import com.ifeanyi.nkataandroid.logic.database.room.model.NetworkBaseUrl

@Database(entities = [Profile::class, Friendship::class, Chat::class, Notification::class, Message::class, Token::class ,NetworkBaseUrl::class], version = 1, exportSchema = false)
abstract class NkataDatabase: RoomDatabase() {

    abstract fun dao(): NkataDatabaseDao

}