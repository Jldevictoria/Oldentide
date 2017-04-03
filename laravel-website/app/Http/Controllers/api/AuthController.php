<?php

namespace App\Http\Controllers\api;

use Illuminate\Http\Request;
use App\Http\Controllers\Controller;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Input;

class AuthController extends Controller
{
	public function login()
	{
		if(Auth::attempt(Input::only('email', 'password'))) {
			return Auth::user();
		}

		return 'invalid username or password';
    }

	public function logout()
	{
		Auth::logout();

		return 'you are logged out';
    }
}
