package com.ifeanyi.nkataandroid.logic.database.room

import androidx.room.TypeConverter
import com.google.gson.Gson
import com.ifeanyi.nkataandroid.logic.database.room.model.MMedia

class ConverterFactory {

    private val gson = Gson()

    @TypeConverter
    fun fromMMedia(media: MMedia?): String? = media?.let { gson.toJson(it) }

    @TypeConverter
    fun toMMedia(json: String?): MMedia? =
        json?.let { gson.fromJson(it, MMedia::class.java) }

}