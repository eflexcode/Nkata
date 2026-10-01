package com.ifeanyi.nkataandroid.ui.screens.navigation

import androidx.compose.runtime.Composable
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import com.ifeanyi.nkataandroid.logic.database.room.model.Notification
import com.ifeanyi.nkataandroid.ui.screens.AddEmailAndVerifyEmailScreen
import com.ifeanyi.nkataandroid.ui.screens.ChatScreen
import com.ifeanyi.nkataandroid.ui.screens.FriendRequestScreen
import com.ifeanyi.nkataandroid.ui.screens.LoginEmailScreen
import com.ifeanyi.nkataandroid.ui.screens.LoginUsernameScreen
import com.ifeanyi.nkataandroid.ui.screens.NotificationScreen
import com.ifeanyi.nkataandroid.ui.screens.ProfileScreen
import com.ifeanyi.nkataandroid.ui.screens.SignUpScreen
import com.ifeanyi.nkataandroid.ui.screens.SplashScreen
import com.ifeanyi.nkataandroid.ui.screens.WelcomeScreen
import com.ifeanyi.nkataandroid.ui.viewmodel.NkataViewModel

@Composable
//fun NtNavGraph(navController: NavHostController, viewModel: NkataViewModel = viewModel()) {
fun NtNavGraph(navController: NavHostController) {

    NavHost(navController, NavScreens.SplashScreenDestination.route) {
        composable(route = NavScreens.SplashScreenDestination.route) {
            SplashScreen()
        }
        composable(route = NavScreens.WelcomeDestination.route) {
            WelcomeScreen()
        }
        composable(route = NavScreens.SignUpScreenDestination.route) {
            SignUpScreen()
        }
        composable(route = NavScreens.LoginUsernameScreenDestination.route) {
            LoginUsernameScreen()
        }
        composable(route = NavScreens.LoginEmailScreenDestination.route) {
            LoginEmailScreen()
        }
        composable(route = NavScreens.AddAndVerifyEmailScreenScreenDestination.route) {
            AddEmailAndVerifyEmailScreen()
        }
        composable(route = NavScreens.ChatScreenDestination.route) {
            ChatScreen()
        }
        composable(route = NavScreens.NotificationScreenDestination.route) {
            NotificationScreen()
        }
        composable(route = NavScreens.ProfileScreenDestination.route) {
            ProfileScreen()
        }
        composable(route = NavScreens.FriendRequestScreenDestination.route) {
            FriendRequestScreen()
        }
        composable(route = NavScreens.FriendRequestScreenDestination.route) {
            FriendRequestScreen()
        }
        composable(route = NavScreens.FriendRequestScreenDestination.route) {
            FriendRequestScreen()
        }
    }

}