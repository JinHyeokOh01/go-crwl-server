CREATE TABLE `cse`
(
    `Number`    VARCHAR(10) NOT NULL COMMENT '게시글 고유 번호',
    `Title`     VARCHAR(200) NOT NULL COMMENT '글 제목',
    `Date`      VARCHAR(20) NOT NULL COMMENT '등록 날짜',
    `Link`      VARCHAR(200) NOT NULL COMMENT '링크',
    PRIMARY KEY (`Number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='컴퓨터공학과 공지사항';

CREATE TABLE `sw`
(
    `Number`    VARCHAR(10) NOT NULL COMMENT '게시글 고유 번호',
    `Title`     VARCHAR(200) NOT NULL COMMENT '글 제목',
    `Date`      VARCHAR(20) NOT NULL COMMENT '등록 날짜',
    `Link`      VARCHAR(200) NOT NULL COMMENT '링크',
    PRIMARY KEY (`Number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='소프트웨어중심대학사업단 공지사항';