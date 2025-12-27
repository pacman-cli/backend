package com.portfolio.blog;

import java.time.LocalDateTime;
import java.util.List;

import org.springframework.boot.CommandLineRunner;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Component;

import com.portfolio.blog.model.BlogPost;
import com.portfolio.blog.model.User;
import com.portfolio.blog.repository.BlogRepository;
import com.portfolio.blog.repository.UserRepository;

@Component
public class DataSeeder implements CommandLineRunner {

  private final UserRepository userRepository;
  private final BlogRepository blogRepository;
  private final PasswordEncoder passwordEncoder;

  public DataSeeder(UserRepository userRepository, BlogRepository blogRepository, PasswordEncoder passwordEncoder) {
    this.userRepository = userRepository;
    this.blogRepository = blogRepository;
    this.passwordEncoder = passwordEncoder;
  }

  @Override
  public void run(String... args) throws Exception {
    // Seed Admin User
    if (userRepository.findByUsername("admin").isEmpty()) {
      User admin = new User();
      admin.setUsername("admin");
      admin.setPassword(passwordEncoder.encode("admin123")); // Default password
      admin.setRole("ROLE_ADMIN");
      admin.setCreatedAt(LocalDateTime.now());
      userRepository.save(admin);
      System.out.println("Admin user seeded: admin / admin123");
    }

    // Seed Blogs
    if (blogRepository.count() == 0) {
      BlogPost post1 = new BlogPost();
      post1.setTitle("Getting Started with Spring Boot and Next.js");
      post1.setSlug("getting-started-spring-boot-nextjs");
      post1.setContent(
          "# Building a Modern Portfolio\n\nThis is a sample blog post demonstrating the power of **Spring Boot** and **Next.js**.\n\n## Why this stack?\n\n1. Type Safety with TypeScript\n2. Robust Backend with Java\n3. Great SEO with Next.js\n\nEnjoy the new site!");
      post1.setTags("java,springboot,nextjs,fullstack");
      post1.setPublished(true);
      post1.setCreatedAt(LocalDateTime.now().minusDays(2));

      BlogPost post2 = new BlogPost();
      post2.setTitle("The Art of Clean Code");
      post2.setSlug("art-of-clean-code");
      post2.setContent(
          "# Clean Code Principles\n\n> \"Clean code always looks like it was written by someone who cares.\"\n\nIn this post, we explore meaningful names, small functions, and S.O.L.I.D principles.");
      post2.setTags("coding,best-practices,clean-code");
      post2.setPublished(true);
      post2.setCreatedAt(LocalDateTime.now().minusDays(1));

      BlogPost post3 = new BlogPost();
      post3.setTitle("Why I Love Framer Motion");
      post3.setSlug("why-i-love-framer-motion");
      post3.setContent(
          "# Animations Made Easy\n\nFramer Motion allows us to create complex animations with declarative syntax.\n\n```jsx\n<motion.div animate={{ x: 100 }} />\n```\n\nIt makes the UI feel alive!");
      post3.setTags("frontend,animation,react");
      post3.setPublished(true);
      post3.setCreatedAt(LocalDateTime.now());

      blogRepository.saveAll(List.of(post1, post2, post3));
      System.out.println("Sample blog posts seeded.");
    }
  }
}
