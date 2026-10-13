package com.ifeanyi.nkataandroid.logic.database.room.model

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "NetworkBaseUrl")
data class NetworkBaseUrl(
    @PrimaryKey(autoGenerate = true)
    val id: Int? = null, val baseUrl: String
)