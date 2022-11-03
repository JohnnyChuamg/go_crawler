# golang

| HttpMethod | Path   | Auth | Memo                        | param                                                      |
| ---- |--------| ---- |-----------------------------|------------------------------------------------------------|
| GET | /ping  | X |                             |                                                            |
| POST | /auth/signIn | X | 登入                          | FROM BODY <br/>{ "username":"johnny","password":"123456" } |
| GET | /v1/ping | X || |
| GET | /v1/work/pic| O/X | 從目標抓一張圖後顯示(不儲存在serversite)  | FROM Query<br/> ?url=https://www.taiwan.net.tw/m1.aspx?sNo=0012076                                      |
|GET | /v1/work/print| O/X | 從目標抓取圖後儲存在server site 在顯示出來 |     FROM Query<br/> ?url=https://www.taiwan.net.tw/m1.aspx?sNo=0012076                                                          |
