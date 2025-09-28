## Backend API

后端API：

### Basic

#### GET `/health`

#### response

```
ok
```

### User

#### POST `/register`

通过用户名和密码注册一个用户

##### request headers

空

##### request body

multipart/form-data:

```
username: a-user-name
password: a-password
```

##### response

``` json
{
  "id": 2,
  "username": "a-user-name",
  "create_at": "2025-09-28T23:41:14.20215884+08:00"
}
```

#### POST `/login`

##### request headers

空

##### request body

multipart/form-data:

```
username: a-user-name
password: a-password
```

##### response

返回一个session_id，用于后续请求的鉴权

``` json
{
  "user_id": 1,
  "token": "TK6PCkWWuF2aptG45o746Lqm3XtHaM_u-UzctSWvvcc"
}
```

#### PUT `/logout`

###### request headers

```
Authorization: your-session-id
```

##### response

空，返回200即为登出成功

### Chat

该端点负责消息的通信

#### GET `/chat/sessions`

##### url query parameters

- num: 查询的条数（最新num条），设置为0则表示请求全部
- from: 一个unix int时间（秒级），表示从该时间点向后寻找

以上两个参数用于分批返回session列表

###### request headers

```
Authorization: your-session-id
```

##### response

按会话更新时间返回该用户（基于认证头）的会话列表

``` json
[
  {
    "session_id": 3,
    "create_at": "2025-09-26T18:02:20+08:00",
    "update_at": "2025-09-28T19:47:29+08:00"
  },
  {
    "session_id": 13,
    "create_at": "2025-09-28T19:35:22+08:00",
    "update_at": "2025-09-28T19:35:27+08:00"
  },
  {
    "session_id": 12,
    "create_at": "2025-09-28T19:23:14+08:00",
    "update_at": "2025-09-28T19:23:36+08:00"
  },
  {
    "session_id": 11,
    "create_at": "2025-09-28T19:22:28+08:00",
    "update_at": "2025-09-28T19:22:28+08:00"
  },
  {
    "session_id": 10,
    "create_at": "2025-09-28T19:18:11+08:00",
    "update_at": "2025-09-28T19:18:11+08:00"
  },
  {
    "session_id": 9,
    "create_at": "2025-09-28T19:17:47+08:00",
    "update_at": "2025-09-28T19:17:47+08:00"
  },
  {
    "session_id": 8,
    "create_at": "2025-09-28T19:17:06+08:00",
    "update_at": "2025-09-28T19:17:06+08:00"
  },
  {
    "session_id": 7,
    "create_at": "2025-09-28T19:12:46+08:00",
    "update_at": "2025-09-28T19:12:46+08:00"
  },
  {
    "session_id": 6,
    "create_at": "2025-09-28T19:07:04+08:00",
    "update_at": "2025-09-28T19:07:04+08:00"
  },
  {
    "session_id": 5,
    "create_at": "2025-09-28T19:06:06+08:00",
    "update_at": "2025-09-28T19:06:06+08:00"
  },
  {
    "session_id": 4,
    "create_at": "2025-09-26T18:04:08+08:00",
    "update_at": "2025-09-26T18:04:08+08:00"
  },
  {
    "session_id": 2,
    "create_at": "2025-09-26T18:00:49+08:00",
    "update_at": "2025-09-26T18:00:49+08:00"
  },
  {
    "session_id": 1,
    "create_at": "2025-09-26T15:47:19+08:00",
    "update_at": "2025-09-26T15:47:19+08:00"
  }
]
```

#### GET `/{session_id}/messages

###### request headers

```
Authorization: your-session-id
```

用户必须是该会话的创建者，否则将返回403

##### response

``` json
[
  {
    "role": "user",
    "content": "The cake is a lie - 2",
    "create_at": "2025-09-28T19:46:34+08:00"
  },
  {
    "role": "assistant",
    "content": "Great! You’ve “entered\" Portal’s most famous Easter‑egg.  \nLet’s unpack what the lie has actually become—and why people love to rattle it with a “‑2”.\n\n---\n\n## 1. The Origin – From Portal to Pop‑Culture\n\n| Year | Medium | How “The cake is a lie” appears |\n|------|--------|---------------------------------|\n| 2007 | **Portal** (video game) | A test chamber’s (Index 2) promise: \"You'll get cake when you finish the tests.\" The cake is a fake reward—players never actually eat it. |\n| 2009 | **Flashing the Gate** (fan‑made short) | Animated meme that exaggerates the lie; the cake is forever out of reach. |\n| 2019 – 2024 | **Social media & memes** | Slides, TikTok dances, Reddit threads, and even deep‑fake animations keep the phrase alive. |\n\n> *“The cake is a lie” has become shorthand for *“this claim is misleading or false”*.  \n> You might see it used by a programmer (“Your API call is a cake is a lie”) or even by a diet‑advisor (“The promised weight‑loss miracle is a cake is a lie”).*\n\n---\n\n## 2. What “‑2” Means (A Few Possible Readings)\n\n| Interpretation | Context | Example |\n|----------------|---------|---------|\n| **Sequel / Second Episode** | In Portal, the original game *”Portal”* followed by *Portal 2* (2007 → 2011). |\n| **Puzzle Hint** | In many puzzle games or escape rooms, “-2” might refer to “remove 2 steps,” implying the lie only holds for the first two layers. |\n| **Community Inside Joke** | Fans sometimes tag “‑2” to denote the second twist: *“Wait, the cake is actually real in level 42!”* |\n| **Metaphorical Layer** | “-2” could suggest the cake lies *below* the surface, like two layers of deception in a political statement. |\n\n---\n\n## 3. A Quick Historical Timeline\n\n```\n2007   – Portal 1 releases (cake promise introduced)\n2009   – First wave of memes (“No cake, just AIs”)\n2011   – Portal 2 where the cake is actually *real* in one scene\n2014   – Microgames (“Portal: How to be a Ghost” uses cake humor)\n2019   – TikTok users create the “Cake Is a Lie” dance\n2022   – Deep‑fake AI demonstrates a *real* cake won’t satisfy\n2024   – Urban myths of “Cake runaway” spawned in gaming forums\n```\n\n---\n\n## 4. Quick “Cake‑Is‑A‑Lie” Quiz (Optional)\n\n1. *What game features a hidden reward that turns out to be a lie?*  \n   **Answer:** Portal\n\n2. *Which character says “We’re in a testing facility.” in Portal 2, hinting at the cake?*  \n   **Answer:** C-3PO (the chicken mannequin) – *just kidding, it’s the “Aperture Science Handheld Portal Device.”*\n\n3. *What is the common tag players add to a meme when they want to *double‑the‑lie*?*  \n   **Answer:** #TheCakeIsALie‑2\n\n---\n\n## 5. Quick React Ideas\n\n| Platform | Meme Format | How to use “The cake is a lie – 2” |\n|-----------|------------|-----------------------------------|\n| **Instagram** | Image of a cracked cake | Caption: *“When you realize the cake is a lie… and then a second cake is found”* |\n| **Twitter** | Quote tweet | Tweet: *“The cake is a lie – 2. Because the first one was a riddle.”* |\n| **YouTube** | Edit a Portal clip | Add a bubble: “You can’t get cake, but you can get ‘cake‑is‑a‑lie‑2’ attempts” |\n\n---\n\n## 6. If You Want Something More\n\n- **A deeper dive** into the lore behind the cake itself?  \n  The cake… (Provides a longer explanation).\n- **A “chef’s special” recipe** to see if you can *actually* make cake?  \n  Sure thing – let’s see if your cake stops lying!\n- **A creative story prompt** where *“The cake is a lie 2”* becomes the twist?  \n  Let me know; I’ll spin something wild.\n\n---\n\n### Bottom Line\n\nThe phrase *“The cake is a lie – 2”* is the community’s version of “Let’s double‑check that.”  \nWhether you’re a Portal veteran, a meme‑hunter, or just feeling nostalgic for a \"gotcha\" line, that lie lives on—often with a clever extra twist.  \n\nSo, ready to see whether any cake is real this time? Let’s roll!",
    "create_at": "2025-09-28T19:47:29+08:00"
  }
]
```

#### POST `/chat/new`

创建一个新回话，发送一条消息，并返回一个llm结果

TODO：注入系统和角色prompt

###### request headers

```
Authorization: your-session-id
```

##### request body

json:

```
{
  "message": {
    "content": "The cake is a lie"
  }
}
```

##### response

``` json
{
  "role": "assistant",
  "content": "You’ve just dropped a classic Portal spoiler! 🎮☕ Whether you’re looking for the secret recipe to the legendary cake or just want to discuss why “the cake is a lie” is the ultimate underdog trope, I’m all ears (and maybe a few virtual crumbs). Let me know what you’d like—file‑limit-eating hints, a nostalgic meme‑deep dive, or a new, legit dessert idea to replace that mythical cake 🍰?",
  "create_at": "2025-09-28T19:35:27.449914424+08:00"
}
```

#### POST `/chat/{session_id}/send`

向已有会话发送新消息，并返回会话id和llm生成结果

TODO：在未生成LLM回复前拒绝用户的send

###### request headers

```
Authorization: your-session-id
```

用户必须是该会话的创建者，否则将返回403

##### request body

{
  "message": {
    "content": "The cake is a lie - 2"
  }
}

##### response

``` json
{
  "session_id": 3,
  "message": {
    "role": "assistant",
    "content": "Great! You’ve “entered\" Portal’s most famous Easter‑egg.  \nLet’s unpack what the lie has actually become—and why people love to rattle it with a “‑2”.\n\n---\n\n## 1. The Origin – From Portal to Pop‑Culture\n\n| Year | Medium | How “The cake is a lie” appears |\n|------|--------|---------------------------------|\n| 2007 | **Portal** (video game) | A test chamber’s (Index 2) promise: \"You'll get cake when you finish the tests.\" The cake is a fake reward—players never actually eat it. |\n| 2009 | **Flashing the Gate** (fan‑made short) | Animated meme that exaggerates the lie; the cake is forever out of reach. |\n| 2019 – 2024 | **Social media & memes** | Slides, TikTok dances, Reddit threads, and even deep‑fake animations keep the phrase alive. |\n\n> *“The cake is a lie” has become shorthand for *“this claim is misleading or false”*.  \n> You might see it used by a programmer (“Your API call is a cake is a lie”) or even by a diet‑advisor (“The promised weight‑loss miracle is a cake is a lie”).*\n\n---\n\n## 2. What “‑2” Means (A Few Possible Readings)\n\n| Interpretation | Context | Example |\n|----------------|---------|---------|\n| **Sequel / Second Episode** | In Portal, the original game *”Portal”* followed by *Portal 2* (2007 → 2011). |\n| **Puzzle Hint** | In many puzzle games or escape rooms, “-2” might refer to “remove 2 steps,” implying the lie only holds for the first two layers. |\n| **Community Inside Joke** | Fans sometimes tag “‑2” to denote the second twist: *“Wait, the cake is actually real in level 42!”* |\n| **Metaphorical Layer** | “-2” could suggest the cake lies *below* the surface, like two layers of deception in a political statement. |\n\n---\n\n## 3. A Quick Historical Timeline\n\n```\n2007   – Portal 1 releases (cake promise introduced)\n2009   – First wave of memes (“No cake, just AIs”)\n2011   – Portal 2 where the cake is actually *real* in one scene\n2014   – Microgames (“Portal: How to be a Ghost” uses cake humor)\n2019   – TikTok users create the “Cake Is a Lie” dance\n2022   – Deep‑fake AI demonstrates a *real* cake won’t satisfy\n2024   – Urban myths of “Cake runaway” spawned in gaming forums\n```\n\n---\n\n## 4. Quick “Cake‑Is‑A‑Lie” Quiz (Optional)\n\n1. *What game features a hidden reward that turns out to be a lie?*  \n   **Answer:** Portal\n\n2. *Which character says “We’re in a testing facility.” in Portal 2, hinting at the cake?*  \n   **Answer:** C-3PO (the chicken mannequin) – *just kidding, it’s the “Aperture Science Handheld Portal Device.”*\n\n3. *What is the common tag players add to a meme when they want to *double‑the‑lie*?*  \n   **Answer:** #TheCakeIsALie‑2\n\n---\n\n## 5. Quick React Ideas\n\n| Platform | Meme Format | How to use “The cake is a lie – 2” |\n|-----------|------------|-----------------------------------|\n| **Instagram** | Image of a cracked cake | Caption: *“When you realize the cake is a lie… and then a second cake is found”* |\n| **Twitter** | Quote tweet | Tweet: *“The cake is a lie – 2. Because the first one was a riddle.”* |\n| **YouTube** | Edit a Portal clip | Add a bubble: “You can’t get cake, but you can get ‘cake‑is‑a‑lie‑2’ attempts” |\n\n---\n\n## 6. If You Want Something More\n\n- **A deeper dive** into the lore behind the cake itself?  \n  The cake… (Provides a longer explanation).\n- **A “chef’s special” recipe** to see if you can *actually* make cake?  \n  Sure thing – let’s see if your cake stops lying!\n- **A creative story prompt** where *“The cake is a lie 2”* becomes the twist?  \n  Let me know; I’ll spin something wild.\n\n---\n\n### Bottom Line\n\nThe phrase *“The cake is a lie – 2”* is the community’s version of “Let’s double‑check that.”  \nWhether you’re a Portal veteran, a meme‑hunter, or just feeling nostalgic for a \"gotcha\" line, that lie lives on—often with a clever extra twist.  \n\nSo, ready to see whether any cake is real this time? Let’s roll!",
    "create_at": "2025-09-28T19:47:29.759618884+08:00"
  }
}
```