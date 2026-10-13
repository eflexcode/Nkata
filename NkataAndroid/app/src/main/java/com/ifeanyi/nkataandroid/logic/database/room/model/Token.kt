package com.ifeanyi.nkataandroid.logic.database.room.model

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "token")
data class Token(@PrimaryKey(autoGenerate = true) val id: Int? = null, val token: String)