package com.ifeanyi.nkataandroid.logic.repository

import com.ifeanyi.nkataandroid.logic.database.room.model.Chat
import com.ifeanyi.nkataandroid.logic.database.room.model.Friendship
import com.ifeanyi.nkataandroid.logic.database.room.model.MMedia
import com.ifeanyi.nkataandroid.logic.database.room.model.Message
import com.ifeanyi.nkataandroid.logic.database.room.model.NetworkBaseUrl
import com.ifeanyi.nkataandroid.logic.database.room.model.Notification
import com.ifeanyi.nkataandroid.logic.database.room.model.Profile
import com.ifeanyi.nkataandroid.logic.database.room.model.Token
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.CheckUsername
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.CheckUsernameResult
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.JwtToken
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.LogInUsername
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.SignUp
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.StandardResponse
import kotlinx.coroutines.flow.Flow

interface LogicRepository {

    fun insertProfile(profile: Profile)

    fun getProfile(): Flow<List<Profile>>

    fun deleteProfile()

//profile end------------------------------------------------------------------------------------------------------------------------------------------------------

    fun insertChat(chat: Chat)

    fun getChat(): Flow<List<Chat>>

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
    )

    fun deleteChat(id: Int)

    fun deleteAllChat()

    //chat end------------------------------------------------------------------------------------------------------------------------------------------------------
    fun insertToken(token: Token)

    fun getToken(): Flow<List<Token>>

    fun deleteToken()

    //token end----------------------------------------------------------------------------------------------------------------------------------------------
    fun insertFriendship(friendship: Friendship)

    fun getFriendships(): Flow<List<Friendship>>

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
    )

    fun deleteFriendship(id: Int)

    fun deleteAllFriendships()

    //friendship end----------------------------------------------------------------------------------------------------------------------------------------------
    fun insertNotification(notification: Notification)

    fun getNotifications(): Flow<List<Notification>>

    fun updateNotificationSeen(
        id: Int,
        seen: Boolean,
        modifiedAt: String
    )

    fun deleteNotification(id: Int)

    //notification end----------------------------------------------------------------------------------------------------------------------------------------------
    fun insertMessage(message: Message)

    fun getMessageBySendUsername(senderUsername: String): Flow<List<Message>>

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
    )

    fun deleteMessage(id: Int)

    fun deleteMessageEmptyChat(senderUsername: String)

    //Message end----------------------------------------------------------------------------------------------------------------------------------------------
    fun insertBaseUrl(baseUrl: NetworkBaseUrl)

    fun getBaseUrl(): Flow<List<NetworkBaseUrl>>

    fun deleteBaseUrl()

    //Local Database end-------------------------------------------------------------------------------------------------------------------------------------------------
    suspend fun signUp(signUp: SignUp): StandardResponse
    suspend fun signInUsername(login: LogInUsername): JwtToken
    suspend fun checkUsername(checkUsername: CheckUsername): CheckUsernameResult
}