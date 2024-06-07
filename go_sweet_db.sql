/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50742 (5.7.42)
 Source Host           : 192.168.253.130:3306
 Source Schema         : go_sweet_db

 Target Server Type    : MySQL
 Target Server Version : 50742 (5.7.42)
 File Encoding         : 65001

 Date: 07/06/2024 23:48:22
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_menu`;
CREATE TABLE `sys_menu`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '採番ID',
  `menu_name` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '菜单名称',
  `router_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '路由名称',
  `menu_type` int(11) NULL DEFAULT NULL COMMENT '权限类型（1：目录；2：菜单）',
  `parent_id` int(11) NULL DEFAULT NULL COMMENT '父菜单ID 若是顶级菜单则为0',
  `is_sys` int(11) NULL DEFAULT NULL COMMENT '是否是系统菜单 1:系统菜单不可删除 0:可删除',
  `order_num` int(11) NULL DEFAULT NULL COMMENT '展示顺序',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'メニュー管理' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_menu
-- ----------------------------
INSERT INTO `sys_menu` VALUES (1, '系统管理', 'system', 1, 0, 1, 99);
INSERT INTO `sys_menu` VALUES (2, '用户管理', 'system_user', 2, 1, 1, 0);
INSERT INTO `sys_menu` VALUES (3, '角色管理', 'system_role', 2, 1, 1, 1);
INSERT INTO `sys_menu` VALUES (4, '权限管理', 'system_permissions', 2, 1, 1, 2);
INSERT INTO `sys_menu` VALUES (7, '首页', 'dashboard', 1, 0, 0, 0);

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '角色名称',
  `deleted` int(11) NULL DEFAULT NULL COMMENT '删除状态（1：删除；０：未删除）',
  `created_by` int(11) NULL DEFAULT NULL,
  `created_date` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `last_modified_by` int(11) NULL DEFAULT NULL,
  `last_modified_date` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'ロール管理' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role
-- ----------------------------

-- ----------------------------
-- Table structure for sys_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_menu`;
CREATE TABLE `sys_role_menu`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NULL DEFAULT NULL COMMENT '角色ID',
  `menu_id` int(11) NULL DEFAULT NULL COMMENT '权限ID',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 23 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '権限管理' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role_menu
-- ----------------------------

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户名',
  `password` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '密码',
  `status` int(11) NULL DEFAULT NULL COMMENT '状态（1：有效；0：无效）',
  `role` int(11) NULL DEFAULT NULL COMMENT '角色ID',
  `deleted` int(11) NULL DEFAULT NULL COMMENT '删除状态（1：删除；０：未删除）',
  `created_by` int(11) NULL DEFAULT NULL COMMENT '创建者',
  `created_date` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '创建时间',
  `last_modified_by` int(11) NULL DEFAULT NULL COMMENT '最后更新者',
  `last_modified_date` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '最后更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'ユーザ管理' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (1, 'root', '03a45ba8294b40fbf58814721b1c21c9', 1, 999999, 0, 1, '2024-06-05 13:24:42', 1, '2024-06-05 13:24:42');

SET FOREIGN_KEY_CHECKS = 1;
