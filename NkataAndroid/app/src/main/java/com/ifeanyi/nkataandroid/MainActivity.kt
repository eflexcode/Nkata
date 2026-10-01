package com.ifeanyi.nkataandroid

import android.Manifest
import android.app.Application
import android.content.Intent
import android.os.Build
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.annotation.RequiresApi
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.slideInVertically
import androidx.compose.animation.slideOutVertically
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.material3.Icon
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.navigation.NavDestination.Companion.hierarchy
import androidx.navigation.NavGraph.Companion.findStartDestination
import androidx.navigation.NavHostController
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import com.ifeanyi.nkataandroid.ui.screens.ProfileScreen
import com.ifeanyi.nkataandroid.ui.screens.navigation.NavScreens
import com.ifeanyi.nkataandroid.ui.screens.navigation.NtNavGraph
import com.ifeanyi.nkataandroid.ui.theme.NkataAndroidTheme
import com.ifeanyi.nkataandroid.ui.utility.UiUtility
import com.ifeanyi.nkataandroid.ui.theme.NavBarLabelColorSelected
import com.ifeanyi.nkataandroid.ui.theme.NavBarLabelColor

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        if (intent?.action == Intent.ACTION_SEND) {
            if ("text/plain" == intent.type) {
                //for text forwarding in
            } else if (intent.type?.startsWith("image/") == true) {
                //for img forwarding  in
            }
        }
        enableEdgeToEdge()
        setContent {
            NkataAndroidTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    Greeting(
                        name = "Android",
                        modifier = Modifier.padding(innerPadding)
                    )
                    ProfileScreen()
                }
            }
        }
    }
}

@Composable
fun Greeting(name: String, modifier: Modifier = Modifier) {
    Text(
        text = "Hello jjj $name!",
        modifier = modifier
    )
}

@RequiresApi(Build.VERSION_CODES.UPSIDE_DOWN_CAKE)
@Composable
fun MainUi(application: Application, mainActivity: MainActivity) {

    val permissionList = arrayOf(
        Manifest.permission.READ_EXTERNAL_STORAGE,
        Manifest.permission.WRITE_EXTERNAL_STORAGE,
        Manifest.permission.INTERNET,
        Manifest.permission.FOREGROUND_SERVICE,
        Manifest.permission.FOREGROUND_SERVICE_LOCATION,
        Manifest.permission.ACCESS_FINE_LOCATION,
        Manifest.permission.ACCESS_COARSE_LOCATION,
        Manifest.permission.RECEIVE_BOOT_COMPLETED,
        Manifest.permission.FOREGROUND_SERVICE_DATA_SYNC,
        Manifest.permission.POST_NOTIFICATIONS,
    )

    val TAG = "MainActivity"
    val navController = rememberNavController()

    Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.BottomCenter) {
        Scaffold(
            bottomBar = { CreateBNV(navController = navController) },
            modifier = Modifier.padding(0.dp)
        ) { inP ->
            Box(modifier = Modifier.padding(inP)) {
                NtNavGraph(navController = navController)
            }
        }
    }

}

@Composable
fun CreateBNV(navController: NavHostController) {

    val navItems = listOf(
        NavScreens.ChatScreenDestination,
        NavScreens.FriendRequestScreenDestination,
        NavScreens.NotificationScreenDestination,
        NavScreens.ProfileScreenDestination
    )

    val navIcons = listOf(
        R.drawable.ic_music,
        R.drawable.ic_album,
        R.drawable.ic_artist,
        R.drawable.ic_list,
    )

    val navIconsSelected = listOf(
        R.drawable.ic_music_selected,
        R.drawable.ic_album_selected,
        R.drawable.ic_artist_selected,
        R.drawable.ic_list_selected
    )

    val navBackStackEntry2 by navController.currentBackStackEntryAsState()
    val currentDestination2 = navBackStackEntry2?.destination

    val isNavBarVisible = currentDestination2?.route != NavScreens.MessageScreenDestination.route || currentDestination2.route != NavScreens.ChatProfileScreenDestination.route

    AnimatedVisibility(
        visible = isNavBarVisible,
        enter = slideInVertically { it },
        exit = slideOutVertically { it }
    ) {
        Column(
            verticalArrangement = Arrangement.spacedBy(0.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            NavigationBar(containerColor = UiUtility.NkataNavBarBackgroundColor, tonalElevation = 40.dp) {

                val navBackStackEntry by navController.currentBackStackEntryAsState()
                val currentDestination = navBackStackEntry?.destination

                navItems.forEachIndexed() { index, items ->

                    var isRout =
                        currentDestination?.hierarchy?.any { items.route == it.route } == true

                    if (items.route == NavScreens.ChatScreenDestination.route &&
                        currentDestination?.route == NavScreens.ChatProfileScreenDestination.route &&
                        currentDestination?.route == NavScreens.MessageScreenDestination.route
                    ) {
                        isRout = true
                    }

                    if (items.route == NavScreens.FriendRequestScreenDestination.route
                    ) {
                        isRout = true
                    }

                    if (items.route == NavScreens.NotificationScreenDestination.route
                    ) {
                        isRout = true
                    }
                    if (items.route == NavScreens.ProfileScreenDestination.route
                    ) {
                        isRout = true
                    }

                    NavigationBarItem(
                        selected = isRout,
                        onClick = {

                            navController.navigate(items.route) {

                                popUpTo(navController.graph.findStartDestination().id) {
                                    saveState = true
                                }

                                launchSingleTop = true
                                restoreState = true

                            }
                        },
                        label = {
                            Text(
                                text = items.title,
                                color = if (isRout) NavBarLabelColorSelected else NavBarLabelColor
                            )
                        },
                        icon = {
                            Icon(
                                painter = painterResource(id = if (isRout) navIconsSelected[index] else navIcons[index]),
                                contentDescription = "", modifier = Modifier
                                    .padding(0.dp)
                                    .size(20.dp),
                                tint = if (isRout) NavBarLabelColorSelected else NavBarLabelColor
                            )
                        }
                    )
                }
            }
        }
    }


}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    NkataAndroidTheme {
        Greeting("Android")
    }
}