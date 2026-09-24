package com.ifeanyi.nkataandroid.logic.network.retrofit.model

import com.google.gson.annotations.SerializedName

data class StandardResponse(@SerializedName("status") val status: String?, @SerializedName("message")val message: String?)