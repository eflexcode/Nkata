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
//                var tokenList =
//                    tokenFlow.asLiveData().value ?: throw NullPointerException("no token found")
                var tokenList =
                    tokenFlow.asLiveData().value

                var token = tokenList?.get(0)

                //retrofit
                val r = Retrofit.Builder()
                    .baseUrl(token?.token)
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
        repo?.insertProfile(profile)
    }

    fun getProfile(): Flow<List<Profile>>? {
        return repo?.getProfile()
    }

    fun deleteProfile() {
        repo?.deleteProfile()
    }

//profile end------------------------------------------------------------------------------------------------------------------------------------------------------

    fun insertChat(chat: Chat) {
        repo?.insertChat(chat)
    }

    fun getChat(): Flow<List<Chat>>? {
        return repo?.getChat()
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
        repo?.updateChat(
            id,
            cloudId,
            displayName,
            email,
            password,
            imgUrl,
            bio,
            isOnline,
            friendsCount,
            groupsCount,
            role,
            createdAt,
            modifiedAt
        )
    }

    fun deleteChat(id: Int) {
        repo?.deleteChat(id)
    }

    fun deleteAllChat() {
        repo?.deleteAllChat()
    }

    //chat end------------------------------------------------------------------------------------------------------------------------------------------------------
    fun insertToken(token: Token) {
        repo?.insertToken(token = token)
    }

    fun getToken(): Flow<List<Token>>? {
        return repo?.getToken()
    }

    fun deleteToken() {
        repo?.deleteToken()
    }

    //token end----------------------------------------------------------------------------------------------------------------------------------------------
    fun insertFriendship(friendship: Friendship) {
        repo?.insertFriendship(friendship)
    }

    fun getFriendships(): Flow<List<Friendship>>? {
        return repo?.getFriendships();
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
        repo?.updateFriendship(
            id,
            cloudId,
            username,
            lastMessage,
            friendUsername,
            friendshipType,
            groupId,
            createdAt,
            modifiedAt
        )
    }

    fun deleteFriendship(id: Int) {
        repo?.deleteFriendship(id)
    }

    fun deleteAllFriendships() {
        repo?.deleteAllFriendships()
    }

    //friendship end----------------------------------------------------------------------------------------------------------------------------------------------
    fun insertNotification(notification: Notification) {
        repo?.insertNotification(notification)
    }

    fun getNotifications(): Flow<List<Notification>>? {
        return repo?.getNotifications()
    }

    fun updateNotificationSeen(
        id: Int,
        seen: Boolean,
        modifiedAt: String
    ) {
        repo?.updateNotificationSeen(id, seen, modifiedAt)
    }

    fun deleteNotification(id: Int) {
        repo?.deleteNotification(id)
    }

    //notification end----------------------------------------------------------------------------------------------------------------------------------------------
    fun insertMessage(message: Message) {
        repo?.insertMessage(message)
    }

    fun getMessageBySendUsername(senderUsername: String): Flow<List<Message>>? {
        return repo?.getMessageBySendUsername(senderUsername)
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
        repo?.updateMessage(
            id,
            cloudId,
            messageID,
            friendshipID,
            senderUsername,
            messageType,
            textContent,
            media,
            createdAt,
            modifiedAt
        )
    }

    fun deleteMessage(id: Int) {
        repo?.deleteMessage(id)
    }

    fun deleteMessageEmptyChat(senderUsername: String) {
        repo?.deleteMessageEmptyChat(senderUsername)
    }

    //Message end----------------------------------------------------------------------------------------------------------------------------------------------
    fun insertBaseUrl(baseUrl: NetworkBaseUrl) {
        repo?.insertBaseUrl(baseUrl)
    }

    fun getBaseUrl(): Flow<List<NetworkBaseUrl>>? {
        return repo?.getBaseUrl()
    }

    fun deleteBaseUrl() {
        repo?.deleteBaseUrl()
    }

    //Local Database end-------------------------------------------------------------------------------------------------------------------------------------------------
    suspend fun signUp(signUp: SignUp): StandardResponse? {
        return repo?.signUp(signUp)
    }

    suspend fun signInUsername(login: LogInUsername): JwtToken? {
        return repo?.signInUsername(login = login)
    }

    suspend fun checkUsername(checkUsername: CheckUsername): CheckUsernameResult? {
        return repo?.checkUsername(checkUsername)
    }

}