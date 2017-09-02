﻿using UnityEngine;
using UnityEngine.SceneManagement;
using System.Collections;
using System.Text.RegularExpressions;
using UnityEngine.UI;
using UnityEngine.EventSystems;

public class AutoLoginUtil : MonoBehaviour
{
    public bool active;
    public string serverIp;
    public string username;
    public string password;

    void Start() {
        if (active) {
            IEnumerator coroutine = SubmitLoginInformation();
            StartCoroutine(coroutine);
        }
    }

    private IEnumerator SubmitLoginInformation()
    {
        string loginUrl = "http://" + serverIp + "/login";

        WWWForm form = new WWWForm();
        form.AddField("login_username", username);
        form.AddField("login_password", password);
        Debug.Log(form.ToString());
        WWW loginResult = new WWW(loginUrl, form);
        yield return loginResult;
        if (!string.IsNullOrEmpty(loginResult.error)) {
            print("Error loging in: " + loginResult.error);
        }
        else {
            Debug.Log(loginResult.text);
            string cookie;
            if (loginResult.responseHeaders.TryGetValue("SET-COOKIE", out cookie))
            {
                Regex sessionFinder = new Regex("session_id=([^;]*)");
                Match sessionId = sessionFinder.Match(cookie);
                Debug.Log(sessionId.Groups[1].Value);
                OldentidePlayerInformation.accountName = username;
                OldentidePlayerInformation.sessionId = sessionId.Groups[1].Value;
            }
            else {
                Debug.Log("User failed to login correctly to the server.");
            }
        }
    }
}
