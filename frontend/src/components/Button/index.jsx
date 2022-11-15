import React from 'react';

import cn from 'classnames';
import style from './Button.module.scss';

const Button = ({
  className = '',
  href = '',
  onClick = () => {},
  children,
  color,
  ...props
}) =>
  href ? (
    <a
      href={href}
      className={cn(style.container, style[color], className)}
      {...props}
    >
      {children}
    </a>
  ) : (
    <button className={cn(style.container, style[color], className)} {...props}>
      {children}
    </button>
  );

export default Button;
