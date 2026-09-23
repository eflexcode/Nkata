package com.ifeanyi.nkataandroid.logic.repository.impl

import com.ifeanyi.nkataandroid.logic.database.room.dao.NkataDatabaseDao
import com.ifeanyi.nkataandroid.logic.database.room.model.Chat
import com.ifeanyi.nkataandroid.logic.database.room.model.Friendship
import com.ifeanyi.nkataandroid.logic.database.room.model.MMedia
import com.ifeanyi.nkataandroid.logic.database.room.model.Message
import com.ifeanyi.nkataandroid.logic.database.room.model.NetworkBaseUrl
import com.ifeanyi.nkataandroid.logic.database.room.model.Notification
import com.ifeanyi.nkataandroid.logic.database.room.model.Profile
import com.ifeanyi.nkataandroid.logic.database.room.model.Token
import com.ifeanyi.nkataandroid.logic.repository.LogicRepository
import kotlinx.coroutines.flow.Flow

class LogicRepositoryImpl(val nkataDatabaseDao: NkataDatabaseDao): LogicRepository {

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
        TODO("Not yet implemented")
    }

    override fun getBaseUrl(): Flow<List<NetworkBaseUrl>> {
        TODO("Not yet implemented")
    }

    override fun deleteBaseUrl() {
        TODO("Not yet implemented")
    }

}