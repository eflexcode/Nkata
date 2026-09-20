package com.ifeanyi.nkataandroid.logic.database.room.model

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "chat")
data class Chat(//that's profile of each firend/friendship
    @PrimaryKey(autoGenerate = true)
    val id: Int? = null,
    val cloudId: Int? = null,
    val displayName: String,
    val email: String,
    val password:  String,
    val imgUrl: String,
    val bio: String,
    val isOnline: String,
    val friendsCount: Int,
    val groupsCount: Int,
    val role: String,
    val createdAt: String,
    val modifiedAt: String
) {
}