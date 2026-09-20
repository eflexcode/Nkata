package com.ifeanyi.nkataandroid.logic.database.room.model

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "friendship")
data class Friendship(
    @PrimaryKey(autoGenerate = true)
    val id: Int? = null,
    val cloudId: Int? = null,
    val username: String,
    val lastMessage: String,
    val friendUsername: String,
    val friendshipType: String,
    val groupId: Int,
    val createdAt: String,
    val modifiedAt: String
)
