package com.ifeanyi.nkataandroid.ui.utility

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color


class UiUtility {
    companion object{
        val NkataButtonColor @Composable get() = if (isSystemInDarkTheme()) Color.Black else Color.White
        val NkataBackgroundColor @Composable get() = if (isSystemInDarkTheme()) Color.Black else Color.White
    }
}