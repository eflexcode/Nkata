package com.ifeanyi.nkataandroid.ui.screens.navigation

sealed class NavScreens(val route: String,
                 val title: String,
                 val iconResourceId: Int) {

    object ChatScreenDestination : NavScreens("chat","Chat",0)
    object FriendRequestScreenDestination : NavScreens("friend_request","FriendRequest",0)
    object NotificationScreenDestination : NavScreens("notification","Notification",0)
    object ProfileScreenDestination : NavScreens("profile","Profile",0)
    object ReceiveShareScreenDestination : NavScreens(" receive_share"," ReceiveShare",0)

}