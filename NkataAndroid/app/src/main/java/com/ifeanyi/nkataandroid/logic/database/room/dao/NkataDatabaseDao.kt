package com.ifeanyi.nkataandroid.logic.database.room.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.Query
import com.ifeanyi.nkataandroid.logic.database.room.model.Chat
import com.ifeanyi.nkataandroid.logic.database.room.model.Friendship
import com.ifeanyi.nkataandroid.logic.database.room.model.MMedia
import com.ifeanyi.nkataandroid.logic.database.room.model.Message
import com.ifeanyi.nkataandroid.logic.database.room.model.NetworkBaseUrl
import com.ifeanyi.nkataandroid.logic.database.room.model.Notification
import com.ifeanyi.nkataandroid.logic.database.room.model.Profile
import com.ifeanyi.nkataandroid.logic.database.room.model.Token
import kotlinx.coroutines.flow.Flow

@Dao
interface NkataDatabaseDao {//can also be called databaseRepository

    @Insert
    fun insertProfile(profile: Profile)

    @Query("SELECT * FROM profile")
    fun getProfile(): Flow<List<Profile>>

    @Query("DELETE FROM profile")
    fun deleteProfile()

//profile end------------------------------------------------------------------------------------------------------------------------------------------------------

    @Insert
    fun insertChat(chat: Chat)

    @Query("SELECT * FROM chat")
    fun getChat(): Flow<List<Chat>>

    @Query("UPDATE chat SET cloudId = :cloudId,displayName=:displayName,email = :email,password =:password,imgUrl=:imgUrl,bio=:bio,isOnline=:isOnline,friendsCount=:friendsCount,groupsCount=:groupsCount,role=:role,createdAt=:createdAt ,modifiedAt=:modifiedAt WHERE id =:id ")
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

    @Query("DELETE FROM chat WHERE id=:id")
    fun deleteChat(id: Int)

    @Query("DELETE FROM chat")
    fun deleteAllChat()

    //chat end------------------------------------------------------------------------------------------------------------------------------------------------------
    @Insert
    fun insertToken(token: Token)

    @Query("SELECT * FROM token")
    fun getToken(): Flow<List<Token>>

    @Query("DELETE FROM token")
    fun deleteToken()

    //token end----------------------------------------------------------------------------------------------------------------------------------------------
    @Insert
    fun insertFriendship(friendship: Friendship)

    @Query("SELECT * FROM friendship")
    fun getFriendships(): Flow<List<Friendship>>

    @Query("UPDATE friendship SET cloudId =:cloudId,username =:username,lastMessage =:lastMessage,friendUsername=:friendUsername,friendshipType=:friendshipType,groupId =:groupId,createdAt=:createdAt,modifiedAt=:modifiedAt WHERE id =:id")
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

    @Query("DELETE FROM friendship WHERE id=:id")
    fun deleteFriendship(id: Int)

    @Query("DELETE FROM friendship")
    fun deleteAllFriendships()

    //friendship end----------------------------------------------------------------------------------------------------------------------------------------------
    @Insert
    fun insertNotification(notification: Notification)

    @Query("SELECT * FROM notification")
    fun getNotifications(): Flow<List<Notification>>

    @Query("UPDATE notification SET seen =:seen, modifiedAt =:modifiedAt WHERE id =:id")
    fun updateNotificationSeen(
        id: Int,
        seen: Boolean,
        modifiedAt: String
    )

    @Query("DELETE FROM notification WHERE id=:id")
    fun deleteNotification(id: Int)

    //notification end----------------------------------------------------------------------------------------------------------------------------------------------
    @Insert
    fun insertMessage(message: Message)

    //TODO update message call on server
    @Query("SELECT * FROM message WHERE senderUsername =:senderUsername ORDER BY createdAt ASC")
    fun getMessageBySendUsername(senderUsername: String): Flow<List<Message>>

    @Query("UPDATE message SET cloudId =:cloudId,messageID =:messageID,friendshipID =:friendshipID,senderUsername=:senderUsername,messageType=:messageType,textContent =:textContent,media=:media,createdAt=:createdAt,modifiedAt=:modifiedAt WHERE id =:id")
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

    @Query("DELETE FROM message WHERE id=:id")
    fun deleteMessage(id: Int)

    @Query("DELETE FROM message WHERE senderUsername=:senderUsername")
    fun deleteMessageEmptyChat(senderUsername: String)

    //Message end----------------------------------------------------------------------------------------------------------------------------------------------
    @Insert
    fun insertBaseUrl(baseUrl: NetworkBaseUrl)

    //TODO update message call on server
    @Query("SELECT * FROM NetworkBaseUrl")
    fun getBaseUrl(): Flow<List<NetworkBaseUrl>>

    @Query("DELETE FROM NetworkBaseUrl")
    fun deleteBaseUrl()

}