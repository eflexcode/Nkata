package com.ifeanyi.nkataandroid.logic.network.retrofit.model

import com.google.gson.annotations.SerializedName

data class JwtToken (@SerializedName("token") val token: String)