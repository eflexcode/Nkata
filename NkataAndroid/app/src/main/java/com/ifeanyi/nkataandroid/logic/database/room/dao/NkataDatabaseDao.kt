package com.ifeanyi.nkataandroid.logic.database.room.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.Query
import com.ifeanyi.nkataandroid.logic.database.room.model.Profile
import kotlinx.coroutines.flow.Flow

@Dao
interface NkataDatabaseDao {//can also be called repository

    @Insert
    suspend fun insertProfile(profile: Profile)

    @Query("SELECT * FROM profile")
    suspend fun getProfile(): Flow<List<Profile>>

    @Query("DELETE FROM profile")
    suspend fun deleteProfile()

//profile end------------------------------------------------------------------------------------------------------------------------------------------------------

    @Insert
    suspend fun insertChat()

}