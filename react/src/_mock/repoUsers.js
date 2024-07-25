import { faker } from '@faker-js/faker';

export const generateMockGithubUsers = (count = 20) => {
  // Randomly decide whether to return users or an empty array
  const shouldReturnUsers = faker.datatype.boolean({ probability: 0.7 }); // 70% chance of returning users

  if (!shouldReturnUsers) {
    return []; // Return an empty array to simulate no users found
  }

  return Array.from({ length: count }, () => {
    const username = faker.internet.userName().toLowerCase();
    return {
      avatar_url: faker.image.avatar(),
      login: username,
      name: faker.person.fullName(),
    };
  });
};
