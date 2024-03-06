import requests

# Function to generate headers with the provided token
def get_headers(token):
    headers = {
        'Authorization': f'Bearer {token}'
    }
    return headers

# Registration data
reg_body = {
    "login": 'adad',
    "email": 'vadimakovlev@ya.ru',
    "password": 'StrongPassword1234',
    "countryCode": 'RU',
    "isPublic": True,
    "phone": "+79999999999",
    'image': ''
}

registerAdmin0 = {
    "login": 'admin8',
    "email": 'admin0@ya.ru',
    "password": 'StrongPassword1234',
    "countryCode": 'RU',
    "isPublic": True,
    "phone": "+7999999",
    'image': 'https://http.cat/images/100.jpg'
}

registervwPP = {
    "login": 'vwPP',
    "email": 'vwPP@ya.ru',
    "password": 'StrongPassword1234',
    "countryCode": 'RU',
    "isPublic": True,
    "phone": "+7999999443",
    'image': 'https://http.cat/images/100.jpg'
}

params = {'limit': 10, 'offset': 0}

# Attempt to register the user
check = requests.post('http://localhost:8080/api/auth/register', json=reg_body)

# Check if registration is successful or user already exists
if check.status_code == 201 or check.status_code == 409:
    print("Registration Successful:", check.text)

    # Login data
    log_body = {
        "login": 'adad',
        "password": 'StrongPassword1234',
    }

    # Attempt to log in
    check2 = requests.post('http://localhost:8080/api/auth/sign-in', json=log_body)

    # Check if login is successful
    if check2.status_code == 200:
        print("Login Successful:", check2.text)
        token_data = check2.json()
        token = token_data.get('token')

        # Fetch user profile
        check3 = requests.get('http://localhost:8080/api/me/profile', headers=get_headers(token))

        # Check if fetching profile is successful
        if check3.status_code == 200:
            print("User Profile:", check3.text)
            
            # Update user profile
            patch_data = {
                "countryCode": "RU",
                "isPublic": True,
                "phone": "+29335528799",
                "image": "https://http.cat/images/100.jpg"
            }
            check4 = requests.patch('http://localhost:8080/api/me/profile', headers=get_headers(token),
                                    json=patch_data)

            # Check if profile update is successful
            if check4.status_code == 200:
                print("Profile Update Successful:", check4.text)

                # Fetch admin's profile
                check5 = requests.get('http://localhost:8080/api/profiles/adad', headers=get_headers(token))

                # Check if fetching admin's profile is successful
                if check5.status_code == 200:
                    print("Admin's Profile:", check5.text)

                    requests.post('http://localhost:8080/api/auth/register', json=registerAdmin0)
                    requests.post('http://localhost:8080/api/auth/register', json=registervwPP)
                    # Continue with subsequent operations

                    # Add Friend
                    addfriend = {'login': 'admin8'}
                    check6 = requests.post('http://localhost:8080/api/friends/add', json=addfriend,
                                           headers=get_headers(token))
                    if check6.status_code == 200:
                        print("Friend Added Successfully:", check6.text)
                    else:
                        print("Failed to add friend:", check6.text)

                    # Add Another Friend
                    addfriend2 = {'login': 'vwPP'}
                    check7 = requests.post('http://localhost:8080/api/friends/add', json=addfriend2,
                                           headers=get_headers(token))
                    if check7.status_code == 200:
                        print("Second Friend Added Successfully:", check7.text)
                    else:
                        print("Failed to add second friend:", check7.text)

                    # Get Friends List
                    check8 = requests.get('http://localhost:8080/api/friends', params=params,
                                           headers=get_headers(token))

                    if check8.status_code == 200:
                        print("Friends List:")
                        friends = check8.json()
                        for friend in friends:
                            print(friend)

                        # Remove Friend
                        check9 = requests.post('http://localhost:8080/api/friends/remove',
                                                json=addfriend2, headers=get_headers(token))
                        if check9.status_code == 200:
                            print("Second Friend Removed Successfully")
                        else:
                            print("Failed to remove second friend:", check9.text)

                        # Publish New Post
                        new_post_data = {"content": "This is a new dadada from admin10", "tags": ["1", "2"]}
                        check10 = requests.post('http://localhost:8080/api/posts/new', json=new_post_data,
                                                headers=get_headers(token))
                        if check10.status_code == 201:
                            print("New Post Published Successfully:", check10.text)
                        else:
                            print("Failed to publish new post:", "Status code is", check10.status_code)

                        # Get Posts by Current User
                        check11 = requests.get('http://localhost:8080/api/posts/feed/my',
                                                headers=get_headers(token))
                        if check11.status_code == 200:
                            print("Posts by Current User:", check11.text)
                        else:
                            print("Failed to fetch posts by current user:", check11.text)

                        # Get Posts by Another User (e.g., admin11)
                        check12 = requests.get('http://localhost:8080/api/posts/feed/admin8',
                                                headers=get_headers(token))
                        if check12.status_code == 200:
                            print("Posts by Another User (admin11==8):", check12.text)
                        else:
                            print("Failed to fetch posts by another user:", check12.text)

                        # Get Post by its ID (let's take the first post for example)
                        first_post_id = check10.json.__get__('id')
                        check13 = requests.get(f'http://localhost:8080/api/posts/{first_post_id}',
                                                headers=get_headers(token), params=params)
                        if check13.status_code == 200:
                            print("Post by ID:", check13.text)
                        else:
                            print("Post by ID:", check13.text)

                        # Like a Post (let's take the first post for example)
                        check14 = requests.post(f'http://localhost:8080/api/posts/{first_post_id}/like',
                                                headers=get_headers(token))
                        if check14.status_code == 200:
                            print("Post Liked:", check14.text)
                        else:
                            print("Post Liked:", check14.text)

                        # Dislike a Post (let's take the first post for example)
                        check15 = requests.post(f'http://localhost:8080/api/posts/{first_post_id}/dislike',
                                                 headers=get_headers(token))
                        if check15.status_code == 200:
                            print("Post Disliked:", check15.text)
                        else:
                            print("Post Disliked:", check14.text)

                    else:
                        print("Failed to get friends list:", check8.text)

                else:
                    print("Failed to fetch admin's profile:", check5.text)

            else:
                print("Failed to update profile:", check4.text)

        else:
            print("Failed to fetch user profile:", check3.text)

    else:
        print("Login failed:", check2.text)

else:
    print("Registration failed:", check.text)