package com.ifeanyi.nkataandroid.ui.screens

import android.window.SplashScreen
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.blur
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import com.ifeanyi.nkataandroid.ui.utility.UiUtility

@Composable
fun SplashScreen(){
    Box(Modifier.background(UiUtility.NkataBackgroundColor).fillMaxSize(), contentAlignment = Alignment.Center){

        Text(text = "Nkata", modifier = Modifier,color= Color.Black, fontWeight = FontWeight(20))

    }
}
