/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // 主色调 - b022hub 暖陶土色
        primary: {
          50: '#fff5ef',
          100: '#ffe7da',
          200: '#ffd0b6',
          300: '#f3ae8a',
          400: '#df8d68',
          500: '#cc785c',
          600: '#b96947',
          700: '#935238',
          800: '#6d3d2b',
          900: '#47261d',
          950: '#2a1510'
        },
        // 辅助色 - 柔和蓝灰，保留一点冷调对比
        accent: {
          50: '#eff5ff',
          100: '#d9e7ff',
          200: '#bdd5ff',
          300: '#93b8f3',
          400: '#6f95d8',
          500: '#5378ba',
          600: '#405f96',
          700: '#304973',
          800: '#223353',
          900: '#151f34',
          950: '#0d1323'
        },
        // 深色模式背景 - 暖黑褐色
        dark: {
          50: '#f8f3e7',
          100: '#efe7d8',
          200: '#ddd2c0',
          300: '#c9beaf',
          400: '#aa9f90',
          500: '#897d70',
          600: '#66594e',
          700: '#3b322c',
          800: '#231e1a',
          900: '#171311',
          950: '#0e0c0b'
        }
      },
      fontFamily: {
        sans: [
          'Manrope',
          'SF Pro Display',
          'PingFang SC',
          'Hiragino Sans GB',
          'Microsoft YaHei',
          'Segoe UI',
          'Helvetica Neue',
          'sans-serif'
        ],
        display: [
          'Cormorant Garamond',
          'Noto Serif SC',
          'Space Grotesk',
          'Manrope',
          'PingFang SC',
          'Hiragino Sans GB',
          'Microsoft YaHei',
          'sans-serif'
        ],
        mono: ['JetBrains Mono', 'SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', 'monospace']
      },
      boxShadow: {
        glass: '0 8px 32px rgba(0, 0, 0, 0.08)',
        'glass-sm': '0 4px 16px rgba(0, 0, 0, 0.06)',
        glow: '0 0 20px rgba(204, 120, 92, 0.28)',
        'glow-lg': '0 0 40px rgba(204, 120, 92, 0.38)',
        card: '0 1px 3px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06)',
        'card-hover': '0 10px 40px rgba(0, 0, 0, 0.08)',
        'inner-glow': 'inset 0 1px 0 rgba(255, 255, 255, 0.1)'
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-primary': 'linear-gradient(135deg, #cc785c 0%, #b96947 100%)',
        'gradient-dark': 'linear-gradient(135deg, #1e293b 0%, #0f172a 100%)',
        'gradient-glass':
          'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
        'mesh-gradient':
          'radial-gradient(at 40% 20%, rgba(204, 120, 92, 0.16) 0px, transparent 50%), radial-gradient(at 80% 0%, rgba(122, 185, 255, 0.08) 0px, transparent 50%), radial-gradient(at 0% 50%, rgba(212, 162, 127, 0.1) 0px, transparent 50%)'
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'slide-down': 'slideDown 0.3s ease-out',
        'slide-in-right': 'slideInRight 0.3s ease-out',
        'scale-in': 'scaleIn 0.2s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        shimmer: 'shimmer 2s linear infinite',
        glow: 'glow 2s ease-in-out infinite alternate'
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' }
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideDown: {
          '0%': { opacity: '0', transform: 'translateY(-10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideInRight: {
          '0%': { opacity: '0', transform: 'translateX(20px)' },
          '100%': { opacity: '1', transform: 'translateX(0)' }
        },
        scaleIn: {
          '0%': { opacity: '0', transform: 'scale(0.95)' },
          '100%': { opacity: '1', transform: 'scale(1)' }
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' }
        },
        glow: {
          '0%': { boxShadow: '0 0 20px rgba(204, 120, 92, 0.28)' },
          '100%': { boxShadow: '0 0 30px rgba(204, 120, 92, 0.42)' }
        }
      },
      backdropBlur: {
        xs: '2px'
      },
      borderRadius: {
        '4xl': '2rem'
      }
    }
  },
  plugins: []
}
