#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

console.log('Installing knas...\n');

// 检查 Go 环境
try {
  execSync('go version', { stdio: 'pipe' });
  console.log('✓ Go environment found');
} catch (e) {
  console.error('✗ Go is not installed');
  console.error('Please install Go from https://golang.org/dl/');
  process.exit(1);
}

// 运行构建
console.log('\nBuilding binaries...');
try {
  execSync('node scripts/build.js', { stdio: 'inherit' });
} catch (e) {
  console.error('✗ Build failed');
  process.exit(1);
}

console.log('\n✓ Installation complete!');
console.log('\nRun "knas init" to configure.');
