-- 创建章节图片表
CREATE TABLE IF NOT EXISTS `anix_chapter_image` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '图片ID',
  `chapter_id` BIGINT NOT NULL COMMENT '章节ID',
  `page_index` INT NOT NULL COMMENT '页码索引，从0开始',
  `image_url` VARCHAR(500) NOT NULL COMMENT '图片URL',
  `width` INT DEFAULT 0 COMMENT '图片宽度',
  `height` INT DEFAULT 0 COMMENT '图片高度',
  `file_size` BIGINT DEFAULT 0 COMMENT '文件大小（字节）',
  `format` VARCHAR(10) DEFAULT 'jpg' COMMENT '图片格式',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节图片表';

-- 创建索引
CREATE INDEX idx_chapter_images_chapter_id ON anix_chapter_image(chapter_id);
CREATE INDEX idx_chapter_images_page_index ON anix_chapter_image(chapter_id, page_index);

-- 插入测试数据
INSERT INTO `anix_chapter_image` (`chapter_id`, `page_index`, `image_url`, `width`, `height`, `file_size`, `format`) VALUES
(1, 0, 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%20action%20scene%20with%20hero%20fighting%20monster&image_size=portrait_4_3', 800, 1200, 256000, 'jpg'),
(1, 1, 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%20dramatic%20dialogue%20scene%20between%20characters&image_size=portrait_4_3', 800, 1200, 245000, 'jpg'),
(1, 2, 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%20emotional%20character%20close%20up%20scene&image_size=portrait_4_3', 800, 1200, 267000, 'jpg'),
(2, 0, 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%20adventure%20scene%20with%20characters%20exploring&image_size=portrait_4_3', 800, 1200, 289000, 'jpg'),
(2, 1, 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%20mystery%20scene%20with%20dark%20atmosphere&image_size=portrait_4_3', 800, 1200, 234000, 'jpg'),
(3, 0, 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%20comedy%20scene%20with%20funny%20characters&image_size=portrait_4_3', 800, 1200, 278000, 'jpg'),
(3, 1, 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%20slice%20of%20life%20peaceful%20scene&image_size=portrait_4_3', 800, 1200, 256000, 'jpg'),
(3, 2, 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%20friendship%20bonding%20scene&image_size=portrait_4_3', 800, 1200, 245000, 'jpg'),
(3, 3, 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%20cliffhanger%20ending%20scene&image_size=portrait_4_3', 800, 1200, 267000, 'jpg');