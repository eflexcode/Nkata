package com.ifeanyi.nkataandroid.ui.screens

import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.lifecycle.viewmodel.compose.viewModel
import com.ifeanyi.nkataandroid.ui.viewmodel.NkataViewModel

@Composable
fun ProfileScreen(viewModel: NkataViewModel = viewModel()) {
  val g = viewModel.g()?.collectAsState(emptyList())
 val f = g?.value
}









