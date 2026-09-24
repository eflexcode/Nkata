package com.ifeanyi.nkataandroid.logic.network.retrofit.model

import com.google.gson.annotations.SerializedName

data class LogInUsername(
    @SerializedName("username") val username: String,
    @SerializedName("password") val password: String
)