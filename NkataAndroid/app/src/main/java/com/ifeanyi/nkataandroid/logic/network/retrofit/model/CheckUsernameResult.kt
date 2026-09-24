package com.ifeanyi.nkataandroid.logic.network.retrofit.model

import com.google.gson.annotations.SerializedName

data class CheckUsernameResult (@SerializedName("exist") val exist: Boolean)