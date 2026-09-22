package com.ifeanyi.nkataandroid.logic.database.room.model

import androidx.room.Entity

@Entity(tableName = "token")
data class Token (val id: Int, val token: String)