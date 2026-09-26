package com.ifeanyi.nkataandroid.logic.network.retrofit.apiservice

import com.ifeanyi.nkataandroid.logic.network.retrofit.model.CheckUsername
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.CheckUsernameResult
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.JwtToken
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.LogInUsername
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.SignUp
import com.ifeanyi.nkataandroid.logic.network.retrofit.model.StandardResponse
import retrofit2.Call
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST

interface NkataRetrofitClient {

    //    r.Route("/auth", func(r chi.Router) {
//        r.Post("/sign-up", apiService.RegisterUser)
//        r.Post("/sign-in-with-username", apiService.SignInUsername)
//        r.Post("/sign-in-with-email", apiService.SignInEmail)
//        r.Post("/sign-in-with-email-verify", apiService.VerifySignInEmailOtp)
//        r.Post("/reset-password", apiService.SendResetPasswordOtp)
//        r.Post("/reset-password-verify", apiService.VerifyResetPasswordOtp)
//        r.Get("/check-username", apiService.CheackUsernameAvailability)
//    })

    //TODO anything with email
    @POST("/auth/sign-up")
    suspend fun signUp(@Body signUp : SignUp): Call<StandardResponse>

    @POST("/auth/sign-in-with-username")
    suspend fun  signInUsername(@Body login: LogInUsername): Call<JwtToken>

    @GET("/auth/check-username")
    suspend fun  checkUsername(@Body checkUsername: CheckUsername): Call<CheckUsernameResult>

}