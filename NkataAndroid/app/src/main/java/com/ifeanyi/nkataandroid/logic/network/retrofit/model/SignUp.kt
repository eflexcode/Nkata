package com.ifeanyi.nkataandroid.logic.network.retrofit.model

import com.google.gson.annotations.SerializedName

data class SignUp(
    @SerializedName("username") val username: String,
    @SerializedName("display_name") val displayName: String,
    @SerializedName("password") val password: String,
)