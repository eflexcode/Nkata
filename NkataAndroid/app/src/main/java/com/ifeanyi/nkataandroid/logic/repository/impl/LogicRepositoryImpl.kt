package com.ifeanyi.nkataandroid.logic.repository.impl

import androidx.compose.ui.Modifier
import com.ifeanyi.nkataandroid.logic.database.room.dao.NkataDatabaseDao
import com.ifeanyi.nkataandroid.logic.database.room.model.Chat
import com.ifeanyi.nkataandroid.logic.database.room.model.Friendship
import com.ifeanyi.nkataandroid.logic.database.room.model.MMedia
import com.ifeanyi.nkataandroid.logic.database.room.model.Message
import com.ifeanyi.nkataandroid.logic.database.room.model.NetworkBaseUrl
import com.ifeanyi.nkataandroid.logic.database.room.model.Notification
import com.ifeanyi.nkataandroid.logic.database.room.model.Profile
import com.ifeanyi.nkataandroid.logic.database.room.model.Token
import com.ifeanyi.nkataandroid.logic.network.retrofit.apiservice.NkataRetrofitClient
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.CheckUsername
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.CheckUsernameResult
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.JwtToken
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.LogInUsername
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.SignUp
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.StandardResponse
import com.ifeanyi.nkataandroid.logic.repository.LogicRepository
import kotlinx.coroutines.flow.Flow
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response
import retrofit2.Retrofit

class LogicRepositoryImpl(val nkataDatabaseDao: NkataDatabaseDao) : LogicRepository {

    lateinit var retrofitClient: NkataRetrofitClient


    override fun insertProfile(profile: Profile) {
        TODO("Not yet implemented")
    }

    override fun getProfile(): Flow<List<Profile>> {
        TODO("Not yet implemented")
    }

    override fun deleteProfile() {
        TODO("Not yet implemented")
    }

    override fun insertChat(chat: Chat) {
        TODO("Not yet implemented")
    }

    override fun getChat(): Flow<List<Chat>> {
        TODO("Not yet implemented")
    }

    override fun updateChat(
        id: Int,
        cloudId: Int,
        displayName: String,
        email: String,
        password: String,
        imgUrl: String,
        bio: String,
        isOnline: String,
        friendsCount: Int,
        groupsCount: Int,
        role: String,
        createdAt: String,
        modifiedAt: String
    ) {
        TODO("Not yet implemented")
    }

    override fun deleteChat(id: Int) {
        TODO("Not yet implemented")
    }

    override fun deleteAllChat() {
        TODO("Not yet implemented")
    }

    override fun insertToken(token: Token) {
        TODO("Not yet implemented")
    }

    override fun getToken(): Flow<List<Token>> {
        TODO("Not yet implemented")
    }

    override fun deleteToken() {
        TODO("Not yet implemented")
    }

    override fun insertFriendship(friendship: Friendship) {
        TODO("Not yet implemented")
    }

    override fun getFriendships(): Flow<List<Friendship>> {
        TODO("Not yet implemented")
    }

    override fun updateFriendship(
        id: Int,
        cloudId: Int,
        username: String,
        lastMessage: String,
        friendUsername: String,
        friendshipType: String,
        groupId: Int,
        createdAt: String,
        modifiedAt: String
    ) {
        TODO("Not yet implemented")
    }

    override fun deleteFriendship(id: Int) {
        TODO("Not yet implemented")
    }

    override fun deleteAllFriendships() {
        TODO("Not yet implemented")
    }

    override fun insertNotification(notification: Notification) {
        TODO("Not yet implemented")
    }

    override fun getNotifications(): Flow<List<Notification>> {
        TODO("Not yet implemented")
    }

    override fun updateNotificationSeen(
        id: Int,
        seen: Boolean,
        modifiedAt: String
    ) {
        TODO("Not yet implemented")
    }

    override fun deleteNotification(id: Int) {
        TODO("Not yet implemented")
    }

    override fun insertMessage(message: Message) {
        TODO("Not yet implemented")
    }

    override fun getMessageBySendUsername(senderUsername: String): Flow<List<Message>> {
        TODO("Not yet implemented")
    }

    override fun updateMessage(
        id: Int?,
        cloudId: Int?,
        messageID: String,
        friendshipID: String,
        senderUsername: String,
        messageType: String,
        textContent: String,
        media: MMedia,
        createdAt: String,
        modifiedAt: String
    ) {
        TODO("Not yet implemented")
    }

    override fun deleteMessage(id: Int) {
        TODO("Not yet implemented")
    }

    override fun deleteMessageEmptyChat(senderUsername: String) {
        TODO("Not yet implemented")
    }

    override fun insertBaseUrl(baseUrl: NetworkBaseUrl) {
        nkataDatabaseDao.insertBaseUrl(baseUrl)
    }

    override fun getBaseUrl(): Flow<List<NetworkBaseUrl>> {
        return nkataDatabaseDao.getBaseUrl()
    }

    override fun deleteBaseUrl() {
        return nkataDatabaseDao.deleteBaseUrl()
    }

    override suspend fun signUp(signUp: SignUp): StandardResponse? {
        var standardResponse: StandardResponse? = StandardResponse("0", "0")
        retrofitClient.signUp(signUp).enqueue(object : Callback<StandardResponse> {
            override fun onResponse(
                call: Call<StandardResponse>,
                response: Response<StandardResponse>
            ) {
                if (response.isSuccessful) {
                    standardResponse = response.body()
                    println(call.request().url.toString() + ": ok")
                } else {
                    //server error
                    standardResponse =
                        StandardResponse(response.code().toString(), response.message())
                    println(call.request().url.toString() + ": server error")
                }
            }

            override fun onFailure(call: Call<StandardResponse?>, t: Throwable) {
                //network error
                standardResponse = StandardResponse("09", "No internet")
                println(call.request().url.toString() + ": network error")

            }
        })
        return standardResponse
    }

    override suspend fun signInUsername(login: LogInUsername): JwtToken? {
        var token: JwtToken? = null
        retrofitClient.signInUsername(login).enqueue(object : Callback<JwtToken> {
            override fun onResponse(call: Call<JwtToken?>, response: Response<JwtToken?>) {
               if (response.isSuccessful){
                   println(call.request().url.toString() + ": ok")
                   token = response.body()

               }else{
                   println(call.request().url.toString() + ": server error")
               }
            }

            override fun onFailure(call: Call<JwtToken?>, t: Throwable) {
                println(call.request().url.toString() + ": network error")
                token = null
            }
        })
        return token
    }

    override suspend fun checkUsername(checkUsername: CheckUsername): CheckUsernameResult? {
        var checkUsernameResult : CheckUsernameResult? = null
        retrofitClient.checkUsername(checkUsername).enqueue(object : Callback<CheckUsernameResult> {
            override fun onResponse(call: Call<CheckUsernameResult?>, response: Response<CheckUsernameResult?>) {
                if (response.isSuccessful){
                    println(call.request().url.toString() + ": ok")
                    checkUsernameResult = response.body()

                }else{
                    println(call.request().url.toString() + ": server error")
                }
            }

            override fun onFailure(call: Call<CheckUsernameResult?>, t: Throwable) {
                println(call.request().url.toString() + ": network error")
                checkUsernameResult = null
            }
        })
        return checkUsernameResult
    }

}