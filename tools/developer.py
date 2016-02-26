import sqlite3

class Developer(object):
    def __init__(self, name, followers, stars, commits, n_repos, avatar):
        self.name = name
        self.followers = followers
        self.stars = stars
        self.commits = commits
        self.n_repos = n_repos
        self.avatar = avatar

    def insert_to_db(self):
        conn = sqlite3.connect(Developer.db_name)
        c = conn.cursor()
        
        c.execute("INSERT INTO developers VALUES ('%s', '%s', '%s', '%s', '%s','%s')"
                    % (self.name, self.followers, self.stars, self.commits, self.n_repos, self.avatar))
        conn.commit()
        conn.close()
        
    @staticmethod
    def init_db(db_name):
        Developer.db_name = db_name
        conn = sqlite3.connect(db_name)
        c = conn.cursor()

        c.execute('''CREATE TABLE developers
                    (name text primary key not null, followers integer, stars integer, commits integer, n_repos integer, avatar text)''')
        conn.commit()
        conn.close()

    def __str__(self):
        return "Developer(%s,%s,%s,%s,%s,%s)" % (self.name, self.followers, self.stars, self.commits, self.n_repos, self.avatar)