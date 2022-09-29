package model

type Book struct {
	Id                 int    `xorm:"not null pk autoincr INT(11)"`
	Title              string `xorm:"not null default '' comment('书籍名称') index VARCHAR(100)"`
	AliasTitle         string `xorm:"not null default '' comment('书籍别名') VARCHAR(255)"`
	StacksBookId       int    `xorm:"not null comment('对应的书库的书籍id') unique INT(10)"`
	SourceId           int    `xorm:"not null default 0 comment('作品所属源站') unique(source_book) INT(10)"`
	ApiBookId          int64  `xorm:"not null default 0 comment('合作方cp的书籍id') unique(source_book) BIGINT(20)"`
	WordsNum           int    `xorm:"not null default 0 comment('书籍总字数') INT(10)"`
	Author             string `xorm:"not null default '' comment('作者名称') VARCHAR(50)"`
	Characters         string `xorm:"not null default '' comment('书籍中主角的名称') VARCHAR(100)"`
	IsOver             int    `xorm:"not null default 0 comment('是否完结（0：未完结    1：已完结）') TINYINT(1)"`
	Category1          int    `xorm:"not null default 0 comment('一级分类id') index(category) INT(10)"`
	Category2          int    `xorm:"not null default 0 comment('二级分类id') index(category) INT(10)"`
	ImageLink          string `xorm:"not null default '' comment('书籍封面') VARCHAR(100)"`
	LatestChapterId    string `xorm:"not null default 0 comment('最新章节的id') INT(10)"`
	LatestChapterTitle string `xorm:"not null default '' comment('最新章节的名称') VARCHAR(100)"`
	LatestChapterUrl   string `xorm:"not null default '' comment('最新章节的跳转地址') VARCHAR(100)"`
	TotalChapterNum    int    `xorm:"not null default 0 comment('总章节数') SMALLINT(6)"`
	OrderNum           int    `xorm:"not null default 0 comment('排序') INT(10)"`
	RankWeek           int    `xorm:"not null default 0 comment('周点击量') INT(10)"`
	RankMonth          int    `xorm:"not null default 0 comment('月点击量') INT(10)"`
	ChapterVer         int    `xorm:"not null default 0 comment('章节的版本') INT(10)"`
	IsBreak            int    `xorm:"not null default 0 comment('是否断更（0：未断更    1：已断更）') TINYINT(1)"`
	IsLock             int    `xorm:"not null comment('是否锁定（0：未锁定    1：人工锁定  2：系统锁定）') TINYINT(1)"`
	IsClassical        int    `xorm:"not null comment('是否经典（0：不是      1：是经典）') TINYINT(1)"`
	IsUp               int    `xorm:"not null default 0 comment('是否上架（0-下架    1-上架）') TINYINT(1)"`
	IsWhite            int    `xorm:"not null default 0 comment('是否白名单（0-不是白名单,1-白名单）') TINYINT(1)"`
	UpdateTime         int    `xorm:"default 0 comment('最后更新时间') INT(11)"`
	Status             int    `xorm:"not null default 0 comment('状态（0:未审核；1：正常；2：回收站）') SMALLINT(3)"`
	//CreatedAt             int     `xorm:"not null default 0 comment('创建时间') INT(10)"`
	//CreatedAt time.Time `xorm:"created"`
	CreatedAt int `xorm:"created"`
	//UpdatedAt             int     `xorm:"not null comment('更新时间') INT(10)"`
	UpdatedAt             int     `xorm:"updated not null comment('更新时间') INT(10)"`
	Level                 int     `xorm:"not null default 0 comment('等级（100-S,90-A,80-B,70-C,50-J,40-L,0-未定义)') SMALLINT(4)"`
	Score                 float64 `xorm:"not null default 7.5 comment('分数') DECIMAL(4,1)"`
	Count                 int     `xorm:"not null default 0 comment('被评为L级的次数') INT(10)"`
	DominantHue           string  `xorm:"not null default '' comment('色调') VARCHAR(20)"`
	ChapterTitleFormat    int     `xorm:"not null default 0 comment('是否格式化章节标题') TINYINT(1)"`
	InitScore             float64 `xorm:"not null default 0.0 comment('书籍脚本评分') DECIMAL(4,1)"`
	AuthorId              int     `xorm:"not null default 0 comment('作者id')  INT(10)"`
	EachAudit             int     `xorm:"not null default 0 comment('是否逐章审核') TINYINT(3)"`
	Price                 float64 `xorm:"not null default 0.0 comment('图书价格') DECIMAL(4,1)"`
	IsHidden              int     `xorm:"not null default 0 comment('是否隐藏（0：未隐藏 1：隐藏)') TINYINT(1)"`
	SubTreasury           int     `xorm:"not null default 4 comment('分库 1:测试库; 2:淘汰库; 3:人工推荐库; 4:系统推荐库') smallint(3)"`
	UpTime                int     `xorm:"default 0 comment('上架时间') INT(11)"`
	CurrentFlowPool       int     `xorm:"default 0 comment('当前流量池') INT(11)"`
	LevelPotential        int     `xorm:"default 0 comment('小说潜力评级') INT(11)"`
	NextChapterUpdateTime string  `xorm:"not null default '' comment('下一章更新时间') VARCHAR(50)"`
}

type BookDetail struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}
