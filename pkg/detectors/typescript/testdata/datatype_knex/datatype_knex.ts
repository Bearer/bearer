import { knex, Knex } from 'knex';
import {User} from './types/user'

knex.select('id').from<User>('users');
knex<User>('users')
  .select('id')
  .select('age')
  .then((users) => {
  });

  declare module 'knex/types/tables' {
    interface User {
      id: number;
      name: string;
      created_at: string;
      updated_at: string;
    }
    
    interface Tables {
      users: User;
      users_composite: Knex.CompositeTableType<
        User,
        Pick<User, 'name'> & Partial<Pick<User, 'created_at' | 'updated_at'>>,
        Partial<Omit<User, 'id'>>
      >;
    }
  }