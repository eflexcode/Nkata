package com.ifeanyi.nkataandroid.logic.database.room.model

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "notification")
data class Notification(
    @PrimaryKey(autoGenerate = true)
    val id: Int? = null,
    val cloudId: Int? = null,
    val userId: String,
    val username: String,
    val tile:  String,
    val message: String,
    val seen: Boolean,
    val createdAt: String,
    val modifiedAt: String
) {
}