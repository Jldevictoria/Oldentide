// Purpose:  Child class representing any user account in Oldentide.

#ifndef OLDENTIDE_ACCOUNT_H_
#define OLDENTIDE_ACCOUNT_H_

#include <string>

class Account {

    protected:

    int id;
    int valid;
    std::string accountName;
    std::string email;
    std::string session;
    int playing;
    std::string key;
    std::string salt;

    public:

    // Constructors.
    Account(int id, int valid, std::string accountName, std::string email,
      std::string session, int playing, std::string key, std::string salt);

    // Getter Functions.
    int GetId();
    int GetValid();
    std::string GetAccountName();
    std::string GetEmail();
    std::string GetSession();
    int GetPlaying();
    std::string GetKey();
    std::string GetSalt();

    // Setter Functions.
    void SetId(int id);
    void SetValid(int valid);
    void SetAccountName(std::string accountName);
    void SetEmail(std::string email);
    void SetSession(std::string session);
    void SetPlaying(int playing);
    void SetKey(std::string key);
    void SetSalt(std::string salt);
};

#endif //OLDENTIDE_ACCOUNT_H_
