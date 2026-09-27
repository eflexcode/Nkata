package com.ifeanyi.nkataandroid.ui.screens.navigation

import androidx.compose.runtime.Composable
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.NavController
import androidx.navigation.NavGraph
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import com.ifeanyi.nkataandroid.ui.viewmodel.NkataViewModel

@Composable
fun NtNavGraph(navController: NavHostController, viewModel: NkataViewModel = viewModel()) {

    NavHost (navController, NavScreens.ChatScreenDestination.route){
        composable(route = NavScreens.ChatScreenDestination.route) {
            ChatScreen()
        }
    }

}