package com.ifeanyi.nkataandroid.logic.database.room.model

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "message")
data class Message(
    @PrimaryKey(autoGenerate = true)
    val id: Int? = null,
    val cloudId: Int? = null,
    val messageID: String,
    val friendshipID: String,
    val senderUsername: String,
    val messageType: String,//MessageChat,MessageReaction,MessageInfo
    val textContent: String,
    val media: MMedia,
    val createdAt: String,
    val modifiedAt: String
) {
}
