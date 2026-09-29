package com.ifeanyi.nkataandroid.ui.screens.navigation

sealed class NavScreens(
    val route: String,
    val title: String,
    val iconResourceId: Int
) {

    object WelcomeDestination : NavScreens(" welcome", " Welcome", 0)
    object SignUpScreenDestination : NavScreens(" signup", " Signup", 0)
    object SplashScreenDestination : NavScreens(" splash", " Splash", 0)
    object LoginEmailScreenDestination : NavScreens(" login_email", " Login Email", 0)
    object LoginUsernameScreenDestination : NavScreens(" login_username", " LoginUsername", 0)
    object AddAndVerifyEmailScreenScreenDestination : NavScreens(" verify", " Verify", 0)

    object ChatScreenDestination : NavScreens("chat", "Chat", 0)
    object FriendRequestScreenDestination : NavScreens("friend_request", "FriendRequest", 0)
    object NotificationScreenDestination : NavScreens("notification", "Notification", 0)
    object ProfileScreenDestination : NavScreens("profile", "Profile", 0)
    object ReceiveShareScreenDestination : NavScreens(" receive_share", " ReceiveShare", 0)

}