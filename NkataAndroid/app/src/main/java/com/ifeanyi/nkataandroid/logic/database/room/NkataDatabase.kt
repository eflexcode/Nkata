package com.ifeanyi.nkataandroid.logic.database.room

import androidx.room.Database
import androidx.room.RoomDatabase
import com.ifeanyi.nkataandroid.logic.database.room.dao.NkataDatabaseDao
import com.ifeanyi.nkataandroid.logic.database.room.model.Chat
import com.ifeanyi.nkataandroid.logic.database.room.model.Friendship
import com.ifeanyi.nkataandroid.logic.database.room.model.Message
import com.ifeanyi.nkataandroid.logic.database.room.model.Notification
import com.ifeanyi.nkataandroid.logic.database.room.model.Profile

@Database(entities = [Profile::class, Friendship::class, Chat::class, Notification::class, Message::class], version = 1)
abstract class NkataDatabase: RoomDatabase() {

    abstract fun dao(): NkataDatabaseDao

}