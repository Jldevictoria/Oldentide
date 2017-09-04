
#include "Account.h"
#include "Utils.h"

//------------------------------------------------------------------------------------------------//
//-------------------                     Class Constructors                   -------------------//
//------------------------------------------------------------------------------------------------//

Account::Account(int id, int valid, std::string accountName, std::string email,
                     std::string session, int playing, std::string key, std::string salt) {
    SetId(id);
    SetValid(valid);
    SetAccountName(accountname);
    SetEmail(email);
    SetSession(session);
    SetPlaying(playing);
    SetKey(key);
    SetSalt(salt);
}

//------------------------------------------------------------------------------------------------//
//-------------------                      Getter Functions                    -------------------//
//------------------------------------------------------------------------------------------------//

std::string Account::GetName() {
    return name;
}

int Account::GetId() {
    return id;
}

int Account::GetValid() {
    return valid;
}

std::string Account::GetAccountName() {
    return accountName;
}

std::string Account::GetEmail() {
    return email;
}

std::string Account::GetSession() {
    return session;
}

int Account::GetPlaying() {
    return playing;
}

std::string Account::GetKey() {
    return key;
}

std::string Account::GetSalt() {
    return salt;
}

//------------------------------------------------------------------------------------------------//
//-------------------                      Setter Functions                    -------------------//
//------------------------------------------------------------------------------------------------//

void Account::SetId(int id) {
    this->id = id;
}

void Account::SetValid(int valid) {
    this->valid = valid;
}

void Account::SetAccountName(std::string accountName) {
    this->accountName = accountName;
}

void Account::SetEmail(std::string email) {
    this->email = email;
}

void Account::SetSession(std::string session) {
    this->session = session;
}

void Account::SetPlaying(int playing) {
    this->playing = playing;
}

void Account::SetKey(std::string key) {
    this->key = key;
}

void Account::SetSalt(std::string salt) {
    this->salt = salt;
}
