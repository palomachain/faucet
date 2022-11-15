import React from 'react';

import cn from 'classnames';
import style from './Section.module.scss';

const Section = ({ children, className = '' }) => (
  <section className={cn(style.container, className)}>{children}</section>
);

export default Section;
