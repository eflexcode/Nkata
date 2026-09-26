package com.ifeanyi.nkataandroid.ui.viewmodel

import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.mutableStateOf
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.asLiveData
import androidx.lifecycle.viewmodel.initializer
import androidx.lifecycle.viewmodel.viewModelFactory
import androidx.room.Room
import com.ifeanyi.nkataandroid.logic.Util
import com.ifeanyi.nkataandroid.logic.database.room.NkataDatabase
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
import com.ifeanyi.nkataandroid.logic.repository.impl.LogicRepositoryImpl
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.collectLatest
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

class NkataViewModel(var repo: LogicRepository?) : ViewModel() {

    private val counter = mutableStateOf(0)

    companion object {
        val Factory = viewModelFactory {
            initializer {
                // Instantiating the repository manually

                val context = this[ViewModelProvider.AndroidViewModelFactory.APPLICATION_KEY]
                    ?: throw IllegalArgumentException("Application context is missing")

                //database
                val appDb = Room.databaseBuilder(
                    context,
                    NkataDatabase::class.java,
                    Util.DatabaseName
                ).build()

                var repo = LogicRepositoryImpl(appDb.dao())

                var tokenFlow = repo.getToken()
                var tokenList =
                    tokenFlow.asLiveData().value ?: throw NullPointerException("no token found")
                var token = tokenList[0]

                //retrofit
                val r = Retrofit.Builder()
                    .baseUrl(token.token)
                    .addConverterFactory(
                        GsonConverterFactory.create()
                    )
                    .build()

                repo.retrofitClient = r.create(NkataRetrofitClient::class.java)
                NkataViewModel(repo)

            }
        }

    }

    fun g(): Flow<List<Token>>? {
        return repo?.getToken()
    }


    fun insertProfile(profile: Profile) {

    }

    fun getProfile(): Flow<List<Profile>> {

    }

    fun deleteProfile() {

    }

//profile end------------------------------------------------------------------------------------------------------------------------------------------------------

    fun insertChat(chat: Chat) {

    }

    fun getChat(): Flow<List<Chat>> {

    }

    fun updateChat(
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

    }

    fun deleteChat(id: Int) {

    }

    fun deleteAllChat() {

    }

    //chat end------------------------------------------------------------------------------------------------------------------------------------------------------
    fun insertToken(token: Token) {
        repo?.insertToken(token = token)
    }

    fun getToken(): Flow<List<Token>> {

    }

    fun deleteToken() {

    }

    //token end----------------------------------------------------------------------------------------------------------------------------------------------
    fun insertFriendship(friendship: Friendship) {

    }

    fun getFriendships(): Flow<List<Friendship>> {

    }

    fun updateFriendship(
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

    }

    fun deleteFriendship(id: Int) {

    }

    fun deleteAllFriendships() {

    }

    //friendship end----------------------------------------------------------------------------------------------------------------------------------------------
    fun insertNotification(notification: Notification) {

    }

    fun getNotifications(): Flow<List<Notification>> {

    }

    fun updateNotificationSeen(
        id: Int,
        seen: Boolean,
        modifiedAt: String
    ) {

    }

    fun deleteNotification(id: Int) {

    }

    //notification end----------------------------------------------------------------------------------------------------------------------------------------------
    fun insertMessage(message: Message) {

    }

    fun getMessageBySendUsername(senderUsername: String): Flow<List<Message>> {

    }

    fun updateMessage(
        id: Int? = null,
        cloudId: Int? = null,
        messageID: String,
        friendshipID: String,
        senderUsername: String,
        messageType: String,//MessageChat,MessageReaction,MessageInfo
        textContent: String,
        media: MMedia,
        createdAt: String,
        modifiedAt: String
    ) {

    }

    fun deleteMessage(id: Int) {

    }

    fun deleteMessageEmptyChat(senderUsername: String) {

    }

    //Message end----------------------------------------------------------------------------------------------------------------------------------------------
    fun insertBaseUrl(baseUrl: NetworkBaseUrl) {

    }

    fun getBaseUrl(): Flow<List<NetworkBaseUrl>> {

    }

    fun deleteBaseUrl() {

    }

    //Local Database end-------------------------------------------------------------------------------------------------------------------------------------------------
    fun signUp(signUp: SignUp): StandardResponse {

    }

    fun signInUsername(login: LogInUsername): JwtToken {

    }

    fun checkUsername(checkUsername: CheckUsername): CheckUsernameResult {

    }

}