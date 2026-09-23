package com.ifeanyi.nkataandroid.logic.database.room.model

import androidx.room.Entity

@Entity(tableName = "NetworkBaseUrl")
data class NetworkBaseUrl (val baseUrl: String)