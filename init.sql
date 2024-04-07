-- This file is part of GoAddons, which is licensed under the GNU General Public License v3.0.
-- You should have received a copy of the GNU General Public License along with this program.
-- If not, see <https://www.gnu.org/licenses/>.

CREATE TABLE IF NOT EXISTS `system_config_default` (
  `name` VARCHAR(255),
  `path` VARCHAR(255),
  PRIMARY KEY (`name`)
);

INSERT OR IGNORE INTO `system_config_default` (`name`, `path`)
  VALUES ('browser.download.dir', '/home/[YOUR_HOST_NAME_HERE]/Downloads/goaddons_download'),
         ('extract.addon.path', '/home/[YOUR_HOST_NAME_HERE]/ssd/Games/battlenet/drive_c/Program Files (x86)/World of Warcraft/_retail_/Interface/AddOns');

CREATE TABLE IF NOT EXISTS `addon` (
  `id` INT AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL UNIQUE,
  `filename` varchar(128) DEFAULT NULL,
  `url` VARCHAR(255) NULL DEFAULT NULL,
  `download_url` VARCHAR(255) NULL DEFAULT NULL,
  `last_downloaded` TIMESTAMP NULL DEFAULT NULL,
  `last_modified_at` TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  `added_at` TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `run_log` (
  `id` INT AUTO_INCREMENT,
  `run_id` VARCHAR(255) NOT NULL,
  `service` VARCHAR(128) NOT NULL,
  `added_at` TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `download_log` (
  `id` INT AUTO_INCREMENT,
  `run_id` VARCHAR(255) NOT NULL,
  `url` VARCHAR(255) NOT NULL,
  `added_at` TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `extract_log` (
  `id` INT AUTO_INCREMENT,
  `run_id` VARCHAR(255) NOT NULL,
  `file` TEXT NOT NULL,
  `added_at` TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY (`id`)
);
